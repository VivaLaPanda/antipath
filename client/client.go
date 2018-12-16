package client

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/VivaLaPanda/antipath/engine"
	"github.com/VivaLaPanda/antipath/engine/action"
	"github.com/VivaLaPanda/antipath/entity"
	"github.com/VivaLaPanda/antipath/entity/player"
	"github.com/VivaLaPanda/antipath/state"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn     *websocket.Conn
	engine   *engine.Engine
	playerID entity.ID

	// Buffered channel of outbound messages.
	stateReciever chan *state.State
}

func NewClient(conn *websocket.Conn, e *engine.Engine) *Client {
	client := &Client{
		conn:          conn,
		engine:        e,
		stateReciever: make(chan *state.State),
	}

	client.playerID = e.AddPlayer()

	return client
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, actionJSON, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Take the message, parse JSON, send as Action to engine
		action := &action.Set{}
		err = json.Unmarshal(actionJSON, action)
		if err != nil {
			log.Printf("error parsing user command: %v", err)
		}

		c.engine.SetAction(c.playerID, *action)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.engine.UnregisterClient(c.playerID)
		ticker.Stop()
		c.conn.Close()
	}()
	// Loop reading current game state
	for {
		select {
		case stateSnapshot, ok := <-c.stateReciever:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The channel is closed.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			clientState := struct {
				ClientData *player.Player
				ClientID   entity.ID
				GameState  *state.State
			}{
				ClientData: c.engine.GetPlayer(c.playerID),
				ClientID:   c.playerID,
				GameState:  stateSnapshot,
			}
			stateString, err := json.Marshal(clientState)
			if err != nil {
				return
			}

			// Need to implement actual JSON encoding logic
			w.Write(stateString)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(e *engine.Engine, w http.ResponseWriter, r *http.Request) {
	log.Println("Client attempting to connect...")
	// No same origin policy
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Make the client
	client := NewClient(conn, e)
	// Register it with the engine
	client.engine.RegisterClient(client.playerID, client.stateReciever)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	log.Printf("Client connected and in game. ID: %v", client.playerID)
}

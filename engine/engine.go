package engine

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/VivaLaPanda/antipath/engine/action"
	"github.com/VivaLaPanda/antipath/entity"
	"github.com/VivaLaPanda/antipath/entity/player"
	"github.com/VivaLaPanda/antipath/state"
)

type Engine struct {
	ClientSubs        map[entity.ID]chan *state.State
	clientSubsLock    *sync.RWMutex
	players           map[entity.ID]*player.Player
	playersLock       *sync.RWMutex
	playerActions     map[entity.ID]action.Set
	playerActionsLock *sync.RWMutex
	actionsToProcess  map[entity.ID]action.Set
	gameState         *state.State
	WindowSize        int
}

func NewEngine(stateSize int, WindowSize int) *Engine {
	engine := &Engine{
		ClientSubs:        make(map[entity.ID]chan *state.State),
		clientSubsLock:    &sync.RWMutex{},
		players:           make(map[entity.ID]*player.Player),
		playersLock:       &sync.RWMutex{},
		playerActions:     make(map[entity.ID]action.Set),
		playerActionsLock: &sync.RWMutex{},
		actionsToProcess:  make(map[entity.ID]action.Set),
		gameState:         state.NewState(stateSize),
		WindowSize:        WindowSize,
	}

	go engine.processEvents()

	return engine
}

func (e *Engine) AddPlayer() (entityID entity.ID) {
	newPlayer := player.NewPlayer()

	// Keep trying to spawn in the player at new coords until it works
	pos := state.Coordinates{
		X: rand.Intn(e.gameState.Size()),
		Y: rand.Intn(e.gameState.Size()),
	}
	entityID, err := e.gameState.NewEntity(newPlayer, pos)
	for err != nil {
		pos := state.Coordinates{
			X: rand.Intn(e.gameState.Size()),
			Y: rand.Intn(e.gameState.Size()),
		}
		entityID, err = e.gameState.NewEntity(newPlayer, pos)
	}

	e.playersLock.Lock()
	newPlayer.PlayerID = entityID
	e.players[entityID] = newPlayer
	e.playersLock.Unlock()
	// Set the default action
	e.playerActionsLock.Lock()
	e.playerActions[entityID] = action.Set{Movement: pos, Jump: false}
	e.playerActionsLock.Unlock()

	return entityID
}

func (e *Engine) RegisterClient(entityID entity.ID, stateReciever chan *state.State) {
	e.clientSubsLock.Lock()
	defer e.clientSubsLock.Unlock()

	e.ClientSubs[entityID] = stateReciever
}

func (e *Engine) UnregisterClient(entityID entity.ID) {
	e.clientSubsLock.Lock()
	defer e.clientSubsLock.Unlock()

	channel := e.ClientSubs[entityID]
	delete(e.ClientSubs, entityID)
	close(channel)
}

func (e *Engine) SetAction(entityID entity.ID, actionSet action.Set) {
	e.playerActionsLock.Lock()
	e.playerActions[entityID] = actionSet
	e.playerActionsLock.Unlock()
}

func (e *Engine) processEvents() {
	for {

		e.processPlayerActions()
		e.updateClients()

		time.Sleep(1000 * time.Millisecond)
	}
}

func (e *Engine) GetPlayer(entityID entity.ID) *player.Player {
	return e.players[entityID]
}

func (e *Engine) processPlayerActions() {
	// Freeze actions at this time
	e.actionsToProcess = e.playerActions
	// Wipe old actions so they don't get reused
	e.playerActionsLock.Lock()
	e.playerActions = make(map[entity.ID]action.Set)
	e.playerActionsLock.Unlock()

	for entityID, action := range e.actionsToProcess {
		var err error
		e.playersLock.RLock()
		playerData := e.players[entityID]
		e.playersLock.RUnlock()

		// Process jumps
		if action.Jump {
			playerData.Jump()
		}

		// Process movement
		// TODO: Enforce player speed
		pos, _ := e.gameState.GetEntityPos(entityID)
		if state.Distance(pos, action.Movement) <= playerData.Speed()*4 {
			err = e.gameState.ChangePos(entityID, action.Movement, playerData.Altitude)
		} else {
			log.Printf("Client %s is moving too fast!", entityID)
		}

		// MovRight now any error is a panic. Once we get to this part of the code actions
		if err != nil {
			panic(err)
		}
	}
}

func (e *Engine) updateClients() {
	e.clientSubsLock.RLock()
	defer e.clientSubsLock.RUnlock()
	for playerID, client := range e.ClientSubs {
		select {
		case client <- e.gameState.PeekState(playerID, e.WindowSize):
		default:
		}
	}
}

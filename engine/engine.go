package engine

import (
	"math/rand"
	"sync"
	"time"

	"github.com/VivaLaPanda/antipath/engine/action"
	"github.com/VivaLaPanda/antipath/entity/player"
	"github.com/VivaLaPanda/antipath/state"
)

type Engine struct {
	players          map[state.EntityID]*player.Player
	ClientSubs       map[state.EntityID]chan *state.State
	clientSubsLock   *sync.RWMutex
	playerActions    map[state.EntityID]action.Set
	actionsToProcess map[state.EntityID]action.Set
	gameState        *state.State
	WindowSize       int
}

func NewEngine(stateSize int, WindowSize int) *Engine {
	engine := &Engine{
		players:          make(map[state.EntityID]*player.Player),
		ClientSubs:       make(map[state.EntityID]chan *state.State),
		clientSubsLock:   &sync.RWMutex{},
		playerActions:    make(map[state.EntityID]action.Set),
		actionsToProcess: make(map[state.EntityID]action.Set),
		gameState:        state.NewState(stateSize),
		WindowSize:       WindowSize,
	}

	go engine.processEvents()

	return engine
}

func (e *Engine) AddPlayer() (entityID state.EntityID) {
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

	e.players[entityID] = newPlayer
	// Set the default action
	e.playerActions[entityID] = action.Set{pos, false}

	return entityID
}

func (e *Engine) RegisterClient(entityID state.EntityID, stateReciever chan *state.State) {
	e.clientSubsLock.Lock()
	defer e.clientSubsLock.Unlock()

	e.ClientSubs[entityID] = stateReciever
}

func (e *Engine) UnregisterClient(entityID state.EntityID) {
	e.clientSubsLock.Lock()
	defer e.clientSubsLock.Unlock()

	channel := e.ClientSubs[entityID]
	delete(e.ClientSubs, entityID)
	close(channel)
}

func (e *Engine) SetAction(entityID state.EntityID, actionSet action.Set) {
	e.playerActions[entityID] = actionSet
}

func (e *Engine) processEvents() {
	for {

		e.processPlayerActions()
		e.updateClients()

		time.Sleep(1000 * time.Millisecond)
	}
}

func (e *Engine) GetPlayer(entityID state.EntityID) *player.Player {
	return e.players[entityID]
}

func (e *Engine) processPlayerActions() {
	// Freeze actions at this time
	e.actionsToProcess = e.playerActions
	// Wipe old actions so they don't get reused
	e.playerActions = make(map[state.EntityID]action.Set)

	for entityID, action := range e.actionsToProcess {
		playerData := e.players[entityID]

		// Process jumps
		if action.Jump {
			playerData.Jump()
		}

		// Process movement
		// TODO: Enforce player speed
		err := e.gameState.ChangePos(entityID, action.Movement, playerData.Altitude)

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

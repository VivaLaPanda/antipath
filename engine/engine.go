package engine

import (
	"math/rand"

	"github.com/VivaLaPanda/antipath/engine/action"
	"github.com/VivaLaPanda/antipath/entity/player"
	"github.com/VivaLaPanda/antipath/state"
	"github.com/VivaLaPanda/antipath/state/tile"
)

type Engine struct {
	players          map[state.EntityID]*player.Player
	ClientSubs       map[chan [][]tile.Tile]bool
	playerActions    map[state.EntityID]action.Set
	actionsToProcess map[state.EntityID]action.Set
	gameState        *state.State
}

func NewEngine() *Engine {
	engine := &Engine{
		players:          make(map[state.EntityID]*player.Player),
		ClientSubs:       make(map[chan [][]tile.Tile]bool),
		playerActions:    make(map[state.EntityID]action.Set),
		actionsToProcess: make(map[state.EntityID]action.Set),
		gameState:        state.NewState(100),
	}

	go engine.processEvents()

	return engine
}

func (e *Engine) AddPlayer() (entityID state.EntityID) {
	newPlayer := player.NewPlayer()
	// Keep trying to spawn in the player at new coords until it works
	var err error
	for err != nil {
		pos := state.Coordinates{
			X: rand.Intn(e.gameState.Size()),
			Y: rand.Intn(e.gameState.Size()),
		}
		entityID, err = e.gameState.NewEntity(newPlayer, pos)
	}

	e.players[entityID] = newPlayer
	// Set the default action
	e.playerActions[entityID] = action.DefaultSet()

	return entityID
}

func (e *Engine) SetAction(entityID state.EntityID, actionSet action.Set) {
	e.playerActions[entityID] = actionSet
}

func (e *Engine) processEvents() {
	for {

		e.processPlayerActions()
	}
}

func (e *Engine) processPlayerActions() {
	// Copy the actions. This enforces the sync and means that if no action
	// is received for the player we will just repeat their last action
	for k, v := range e.playerActions {
		e.actionsToProcess[k] = v
	}
	for entityID, action := range e.playerActions {
		playerData := e.players[entityID]

		// Process jumps
		if action.Jump {
			playerData.Jump()
		}

		// Process movement
		err := e.gameState.Move(entityID, action.Movement, playerData.Speed(), playerData.Altitude)

		// Right now any error is a panic. Once we get to this part of the code actions
		if err != nil {
			panic(err)
		}
	}
}

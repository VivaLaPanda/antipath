package engine

import (
	"math/rand"

	"github.com/VivaLaPanda/antipath/engine/action"
	"github.com/VivaLaPanda/antipath/entity/player"
	"github.com/VivaLaPanda/antipath/grid"
)

type Engine struct {
	players          map[string]*player.Player
	playerActions    map[string]action.Set
	actionsToProcess map[string]action.Set
	gameState        *grid.State
}

func NewEngine() *Engine {
	engine := &Engine{
		playerActions:    make(map[string]action.Set),
		actionsToProcess: make(map[string]action.Set),
		gameState:        grid.NewState(100),
	}

	go engine.processEvents()

	return engine
}

func (e *Engine) AddPlayer() (playerID string) {
	newPlayer := player.NewPlayer()
	// Keep trying to spawn in the player at new coords until it works
	var err error
	for err != nil {
		pos := grid.Coordinates{
			X: rand.Intn(e.gameState.Size()),
			Y: rand.Intn(e.gameState.Size()),
		}
		playerID, err = e.gameState.NewEntity(newPlayer, pos)
	}

	e.players[playerID] = newPlayer
	// Set the default action
	e.playerActions[playerID] = action.DefaultSet()

	return playerID
}

func (e *Engine) SetAction(playerID string, actionSet action.Set) {
	e.playerActions[playerID] = actionSet
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
	for playerID, action := range e.playerActions {
		playerData := e.players[playerID]

		// Process jumps
		if action.Jump {
			playerData.Jump()
		}

		// Process movement
		err := e.gameState.Move(playerID, action.Movement, playerData.Speed(), playerData.Altitude)

		// Right now any error is a panic. Once we get to this part of the code actions
		if err != nil {
			panic(err)
		}
	}
}

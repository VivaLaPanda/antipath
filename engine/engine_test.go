package engine

import (
	"testing"
	"time"

	"github.com/VivaLaPanda/antipath/engine/action"
	"github.com/VivaLaPanda/antipath/state"
)

func TestNewEngine(t *testing.T) {
	_ = NewEngine(100, 20)
	return
}

func TestAddPlayer(t *testing.T) {
	engine := NewEngine(100, 20)
	id := engine.AddPlayer()

	if engine.players[id].Health != 100 {
		t.Errorf("Newly added player has non-default values. %v", engine.players[id])
	}

	_, exists := engine.gameState.GetEntityPos(id)
	if !exists {
		t.Errorf("Newly added entity not added to the game state.")
	}

	return
}

func TestSetAction(t *testing.T) {
	engine := NewEngine(100, 20)
	id := engine.AddPlayer()

	action := action.DefaultSet()
	engine.SetAction(id, action)

	action.Movement = state.MovUp

	for idx := 0; idx < 50; idx++ {
		engine.gameState.GetEntityPos(id)
		engine.SetAction(id, action)
	}

	return
}

func TestClientSubs(t *testing.T) {
	engine := NewEngine(50, 10)
	id := engine.AddPlayer()
	for idx := 0; idx < 20; idx++ {
		engine.AddPlayer()
	}
	stateReciever := make(chan *state.State)

	// Add the subscription
	engine.RegisterClient(id, stateReciever)

	// Watch the state updates
	go func() {
		for {
			_, ok := <-stateReciever
			if !ok {
				break
			}
		}
	}()

	// Just move up a bunch
	action := action.DefaultSet()
	action.Movement = state.MovUp

	for idx := 0; idx < 25; idx++ {
		engine.SetAction(id, action)
		time.Sleep(10 * time.Millisecond)
	}

	engine.UnregisterClient(id)

	return
}

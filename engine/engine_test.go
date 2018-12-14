package engine

import (
	"testing"

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

	pos, _ := engine.gameState.GetEntityPos(id)

	testAction := action.Set{pos, false}
	engine.SetAction(id, testAction)

	for idx := 0; idx < 50; idx++ {
		pos.Y -= 1
		testAction = action.Set{pos, false}
		engine.SetAction(id, testAction)
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

	pos, _ := engine.gameState.GetEntityPos(id)

	testAction := action.Set{Movement: pos, Jump: false}
	engine.SetAction(id, testAction)

	for idx := 0; idx < 50; idx++ {
		pos.Y -= 1
		testAction = action.Set{Movement: pos, Jump: false}
		engine.SetAction(id, testAction)
	}

	engine.UnregisterClient(id)

	return
}

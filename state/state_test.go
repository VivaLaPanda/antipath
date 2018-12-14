package state

import (
	"testing"

	"github.com/VivaLaPanda/antipath/entity/player"
)

func TestNewState(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic while checking if new state is valid: Err: %v", r)
		}
	}()

	NewState(100)
}

func TestMarshallJson(t *testing.T) {
	testState := NewState(100)

	_, err := testState.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal player struct into JSON, err: %v", err)
	}
}

func TestGetTile(t *testing.T) {
	testState := NewState(100)

	testPos := Coordinates{0, 0}
	tile, err := testState.GetTile(testPos)
	if err != nil {
		t.Errorf("Getting valid tile produced err: %v", err)
		return
	}
	if tile.Height() != 0 {
		t.Errorf("Tile not properly initialized")
	}

	testBadPos := Coordinates{200, 200}
	tile, err = testState.GetTile(testBadPos)
	if err == nil {
		t.Errorf("Indexing off of the state should result in an error")
	}
}

func TestNewEntity(t *testing.T) {
	testState := NewState(100)
	pos := Coordinates{0, 0}
	player := player.NewPlayer()
	_, err := testState.NewEntity(player, pos)
	if err != nil {
		t.Errorf("Placing a valid entity into the state at a valid pos produced an error: %v", err)
	}
	_, err = testState.NewEntity(player, pos)
	if err == nil {
		t.Errorf("Placing an entity into a full space didn't produce an error")
	}
}

func TestGetEntityPos(t *testing.T) {
	testState := NewState(100)
	pos := Coordinates{0, 0}
	player := player.NewPlayer()
	playerID, _ := testState.NewEntity(player, pos)
	actualPos, exists := testState.GetEntityPos(playerID)
	if exists == false || pos != actualPos {
		t.Errorf("Was unable to properly fetch newly created entity's position")
	}
}

func TestPeekState(t *testing.T) {
	testState := NewState(100)
	pos := Coordinates{0, 0}
	player := player.NewPlayer()
	playerID, _ := testState.NewEntity(player, pos)

	stateFrag := testState.PeekState(playerID, 10)
	if len(stateFrag.grid) > 10 {
		t.Errorf("Peekstate resulted in bigger grid than expected. A: %d, E:%d", len(stateFrag.grid), 10)
	}

	actualPos, _ := stateFrag.entities[playerID]
	if actualPos != pos {
		t.Errorf("Peekstate result doesn't have correct entity pos. A:%v E:%v", actualPos, pos)
	}
}

func TestMove(t *testing.T) {
	testState := NewState(100)
	pos := Coordinates{50, 50}
	testPlayer := player.NewPlayer()
	playerID, err := testState.NewEntity(testPlayer, pos)

	// Check a small move
	err = testState.Move(playerID, MovUp, testPlayer.Speed(), testPlayer.Altitude)
	if err != nil {
		t.Errorf("Moving the player to an empty space resulted in an error")
	}
	newPos, _ := testState.GetEntityPos(playerID)
	expectedPos := pos
	expectedPos.Y -= testPlayer.Speed()
	if newPos != expectedPos {
		t.Errorf("Move didn't result in the expected location. A: %v, E: %v", newPos, expectedPos)
	}

	// Try and hit the top wall
	err = testState.Move(playerID, MovUp, 100, testPlayer.Altitude)
	if err != nil {
		t.Errorf("Moving the player resulted in an error")
	}
	newPos, _ = testState.GetEntityPos(playerID)
	expectedPos.Y = 0
	if newPos != expectedPos {
		t.Errorf("Move didn't result in the expected location. A: %v, E: %v", newPos, expectedPos)
	}

	// Test collision against other player
	pos = Coordinates{55, 0}
	otherPlayer := player.NewPlayer()
	testState.NewEntity(otherPlayer, pos)

	err = testState.Move(playerID, MovRight, 10, testPlayer.Altitude)
	if err != nil {
		t.Errorf("Moving the player resulted in an error")
	}
	newPos, _ = testState.GetEntityPos(playerID)
	expectedPos.X = 54
	if newPos != expectedPos {
		t.Errorf("Move didn't result in the expected location. A: %v, E: %v", newPos, expectedPos)
	}
}

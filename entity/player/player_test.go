package player

import (
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	testPlayer := NewPlayer()

	_, err := testPlayer.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal player struct into JSON, err: %v", err)
	}
}

func TestJump(t *testing.T) {
	testPlayer := NewPlayer()

	testPlayer.Jump()
	if testPlayer.Altitude != (testPlayer.jumpHeight + 1) {
		t.Errorf("Player jump resulted in wrong altitude")
	}
}

func TestFall(t *testing.T) {
	testPlayer := NewPlayer()

	testPlayer.Jump()
	testPlayer.Fall(1)

	if testPlayer.Altitude != testPlayer.jumpHeight {
		t.Errorf("Player didn't fall at the expected speed")
	}
}

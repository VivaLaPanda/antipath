package tile

import (
	"testing"

	"github.com/VivaLaPanda/antipath/entity/player"
)

func TestSetEntity(t *testing.T) {
	testTile := Tile{}
	testEntity := player.NewPlayer()

	err := testTile.SetEntity(testEntity)
	if err != nil {
		t.Errorf("Error setting entity on empty tile")
		return
	}

	err = testTile.SetEntity(testEntity)
	if err == nil {
		t.Errorf("No error setting entity on tile that already has one")
	}
}

func TestPopEntity(t *testing.T) {
	testTile := Tile{}
	testEntity := player.NewPlayer()

	err := testTile.SetEntity(testEntity)
	if err != nil {
		t.Errorf("Error setting entity on empty tile")
		return
	}

	resultEntity := testTile.PopEntity()
	if resultEntity != testEntity {
		t.Errorf("Push and then pop don't match.")
		return
	}

	resultEntity = testTile.PopEntity()
	if resultEntity != nil {
		t.Errorf("Pop on an empty tile isn't nil.")
		return
	}

	err = testTile.SetEntity(testEntity)
	if err != nil {
		t.Errorf("Error setting entity when it should be empty")
		return
	}
}

func TestWillCollide(t *testing.T) {
	testTileA := Tile{height: 0}
	testTileB := Tile{height: 100}
	if testTileA.WillCollide(50) == true {
		t.Errorf("a 50 elevation enitity collided with a 0 height tile.")
	}
	if testTileA.WillCollide(0) == false {
		t.Errorf("a 0 elevation enitity didn't collide with a 0 height tile.")
	}
	if testTileB.WillCollide(150) == true {
		t.Errorf("a 150 elevation enitity collided with a 100 height tile.")
	}
	if testTileB.WillCollide(50) == false {
		t.Errorf("a 50 elevation enitity didn't collide with a 100 height tile.")
	}

	testEntity := player.NewPlayer()

	testTileA.entity = testEntity

	if testTileA.WillCollide(5) == false {
		t.Errorf("A 0 height tile with a 10 height player should collide with a 5 elevation object")
	}
	if testTileA.WillCollide(15) == true {
		t.Errorf("A 0 height tile with a 10 height player should not collide with a 15 elevation object")
	}
}

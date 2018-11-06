package grid

import (
	"fmt"

	"github.com/VivaLaPanda/antipath/entity"
	"github.com/VivaLaPanda/antipath/grid/tile"
	uuid "github.com/satori/go.uuid"
)

type State struct {
	grid     [][]tile.Tile
	size     int
	entities map[string]Coordinates
}

type Coordinates struct {
	X, Y int
}

type Direction int

const (
	Up    Direction = iota
	Right Direction = iota
	Left  Direction = iota
	Down  Direction = iota
)

func NewState(size int) (grid *State) {
	gridData := make([][]tile.Tile, size)
	for idx := range gridData {
		gridData[idx] = make([]tile.Tile, size)
	}
	return &State{
		grid:     gridData,
		size:     size, // faster than using len every time
		entities: make(map[string]Coordinates),
	}
}

func (s *State) GetTile(pos Coordinates) (*tile.Tile, error) {
	if outOfBounds(s.size, pos) {
		return nil, fmt.Errorf("provided pos is out of bounds. Pos: %v, maxsize: %d", pos, s.size)
	}
	return &s.grid[pos.X][pos.Y], nil
}

func (s *State) NewEntity(data entity.Entity, pos Coordinates) (id string, err error) {
	targetTile, err := s.GetTile(pos)
	if err != nil {
		return "", err
	}

	if err := targetTile.SetEntity(data); err != nil {
		return "", fmt.Errorf("provided pos can't contain an entity, already full. Tile %v", targetTile)
	}

	id = uuid.Must(uuid.NewV4()).String()

	s.entities[id] = pos

	return id, nil
}

func (s *State) Move(entityID string, dir Direction, speed int, instant bool) (err error) {
	// Get the location of the entity
	sourcePos, exists := s.entities[entityID]
	if !exists {
		return fmt.Errorf("provided entity ID not valid. ID: %s", entityID)
	}
	// Get the tile data at that location
	sourceTile, err := s.GetTile(sourcePos)
	if err != nil {
		return fmt.Errorf("couldn't get tile at provided pos, pos: %v, err: %s", sourcePos, err)
	}

	// Calculate the movement accounting for instant
	var targetPos Coordinates
	var targetTile *tile.Tile
	for targetPos = sourcePos; speed > 0; speed -= 1 {
		posDelta := 1
		if instant {
			posDelta = speed
		}
		switch dir {
		case Up:
			targetPos.Y += posDelta
		case Down:
			targetPos.Y -= posDelta
		case Left:
			targetPos.X -= posDelta
		case Right:
			targetPos.X += posDelta
		}

		// Get tile data for where we moved to
		targetTile, err = s.GetTile(targetPos)
		if err != nil {
			return fmt.Errorf("couldn't move to resulting pos: %v, err: %s", targetPos, err)
		}

		// Make sure out target is free
		if !targetTile.IsFree() {
			return fmt.Errorf("couldn't move to resulting pos: %v, tile is occupied", targetPos)
		}

		if instant {
			break
		}
	}

	// Move the entity
	entityData := sourceTile.PopEntity()
	targetTile.SetEntity(entityData)
	s.entities[entityID] = targetPos

	return nil
}

func outOfBounds(size int, pos Coordinates) bool {
	return pos.X > size || pos.Y > size
}

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

func (s *State) Move(entityID string, dir Direction, speed int, altitude int) (err error) {
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

	// Calculate the total movement
	var targetPos Coordinates
	var targetTile *tile.Tile
	switch dir {
	case Up:
		targetPos.Y += speed
	case Down:
		targetPos.Y -= speed
	case Left:
		targetPos.X -= speed
	case Right:
		targetPos.X += speed
	}

	// Simulate entity movement with collision rules
	resultPos := s.moveCollider(sourcePos, targetPos, altitude)
	targetTile, _ = s.GetTile(resultPos)

	// Move the entity
	entityData := sourceTile.PopEntity()
	targetTile.SetEntity(entityData)
	s.entities[entityID] = targetPos

	return nil
}

func (s *State) moveCollider(sourcePos Coordinates, targetPos Coordinates, altitude int) (result Coordinates) {
	// Keep track of our movements
	result = sourcePos
	checkPos := sourcePos
	// Loop counter is simply in case some bug causes an infinite loop
	// If anything moves a distance greater than twice the total board size
	// something is wrong
	for distanceMoved := 0; distanceMoved < s.size*2; distanceMoved++ {
		// move 1 towards out destination. If we're already at our destination
		// just return that
		switch {
		case targetPos.X > checkPos.X:
			result.X += 1
		case targetPos.X < checkPos.X:
			result.X -= 1
		case targetPos.Y > checkPos.Y:
			result.Y += 1
		case targetPos.Y < checkPos.Y:
		default: // Positions are the same
			return targetPos
		}

		// Get tile data for where we moved to
		checkTile, err := s.GetTile(checkPos)
		if err != nil {
			return result
		}

		// Make sure out target is free
		if !checkTile.CheckCollision(altitude) {
			return result
		}

		// Store that we successfully can move here
		result = checkPos
	}

	panic("movement calculation out of bounds!")
}

func outOfBounds(size int, pos Coordinates) bool {
	return pos.X > size || pos.Y > size
}

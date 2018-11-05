package grid

import (
	"github.com/VivaLaPanda/antipath/entity"
	"github.com/VivaLaPanda/antipath/grid/tile"
)

type State struct {
	grid     [][]tile.Tile
	entities map[string]entity.Entity
}

func NewState(size int) (grid *State) {
	gridData := make([][]tile.Tile, size)
	for idx := range gridData {
		gridData[idx] = make([]tile.Tile, size)
	}
	return &State{
		grid:     gridData,
		entities: make(map[string]entity.Entity),
	}
}

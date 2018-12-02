package tile

import (
	"fmt"

	"github.com/VivaLaPanda/antipath/entity"
)

type Tile struct {
	alignment int
	entity    entity.Entity
	//hazard *Hazard
	totemHealth    int
	alignmentDelta int
	height         int
}

func (tile *Tile) SetEntity(entity entity.Entity) error {
	if tile.entity != nil {
		return fmt.Errorf("can only SetEntity if entity is already nil, remove before setting")
	}
	tile.entity = entity

	return nil
}

func (tile *Tile) PopEntity() entity.Entity {
	var ref entity.Entity // declare so the pointer logic is a little clearer
	ref = tile.entity
	tile.entity = nil
	return ref
}

func (tile *Tile) Height() int {
	if tile.entity != nil {
		return tile.height + tile.entity.Height()
	}
	return tile.height
}

func (tile *Tile) WillCollide(altitude int) bool {
	return altitude <= tile.Height()
}

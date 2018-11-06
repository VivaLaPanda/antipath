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
}

func (tile *Tile) IsFree() bool {
	if tile.totemHealth > 0 {
		return false
	}
	if tile.entity != nil {
		return false
	}

	return true
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
	return &ref
}

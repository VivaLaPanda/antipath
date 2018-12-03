package action

import "github.com/VivaLaPanda/antipath/grid"

type Set struct {
	Movement grid.Direction
	Jump     bool
}

func DefaultSet() Set {
	return Set{
		Movement: grid.None,
		Jump:     false,
	}
}

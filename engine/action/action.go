package action

import "github.com/VivaLaPanda/antipath/state"

type Set struct {
	Movement state.Direction
	Jump     bool
}

func DefaultSet() Set {
	return Set{
		Movement: state.MovNone,
		Jump:     false,
	}
}

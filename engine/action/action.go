package action

import "github.com/VivaLaPanda/antipath/state"

type Set struct {
	Movement  state.Coordinates
	Jump      bool
	Attack    int
	AttackDir state.Direction
}

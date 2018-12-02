package player

type Player struct {
	Health    uint
	alignment int
	baseSpeed int
	height    int
	Altitude  int
}

func NewPlayer() *Player {
	return &Player{
		Health:    100,
		alignment: 0,
		baseSpeed: 1,
		height:    5,
		Altitude:  1,
	}
}

func (p *Player) Height() int {
	return p.height
}

func (p *Player) Speed() int {
	return p.baseSpeed
}

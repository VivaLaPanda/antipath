package player

type Player struct {
	Health     uint
	alignment  int
	baseSpeed  int
	height     int
	jumpHeight int
	Altitude   int
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

func (p *Player) Jump() {
	// You can only jump if you'r already on the ground
	if p.Altitude == 1 {
		p.Altitude += p.jumpHeight
	}
}

func (p *Player) Fall(speed int) {
	if p.Altitude-speed > 1 {
		p.Altitude -= speed
	} else {
		p.Altitude = 1
	}
}

func (p *Player) Height() int {
	return p.height
}

func (p *Player) Speed() int {
	return p.baseSpeed
}

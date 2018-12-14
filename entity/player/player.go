package player

import "encoding/json"

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
		baseSpeed: 5,
		height:    5,
		Altitude:  1,
	}
}

func (p *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Health     uint `json:"health"`
		Alignment  int  `json:"alignment"`
		Speed      int  `json:"speed"`
		Height     int  `json:"height"`
		JumpHeight int  `json:"jumpHeight"`
		Altitude   int  `json:"altitude"`
	}{
		Health:     p.Health,
		Alignment:  p.alignment,
		Speed:      p.Speed(),
		Height:     p.Height(),
		JumpHeight: p.jumpHeight,
		Altitude:   p.Altitude,
	})
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

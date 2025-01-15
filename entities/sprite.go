package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img       *ebiten.Image
	X         float64
	Y         float64
	VX        float64
	JumpState uint
	IsLeft    bool
}

package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Component interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int)
}

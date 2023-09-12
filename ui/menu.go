package ui

import (
	"image/color"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/hajimehoshi/ebiten/v2"
)

type Menu struct {
	width    int
	height   int
	children []Component
}

func NewMenu(children []Component) *Menu {
	width := constant.MenuWidth
	_, height := ebiten.WindowSize()
	return &Menu{
		width:    width,
		height:   height,
		children: children,
	}
}

func (m *Menu) Update() error {
	_, m.height = ebiten.WindowSize()
	for _, c := range m.children {
		c.Update()
	}
	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(m.width, screen.Bounds().Dy())
	img.Fill(color.Gray{Y: constant.MenuGrayY})
	screen.DrawImage(img, nil)
	for _, c := range m.children {
		c.Draw(screen)
	}
}

func (m *Menu) Layout(outsideWidth, outsideHeight int) {
	for _, c := range m.children {
		c.Layout(outsideWidth, outsideHeight)
	}
}

package main

import (
	"log"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/aethiopicuschan/odori/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(constant.DefaultScreenWidth, constant.DefaultScreenHeight)
	ebiten.SetWindowSizeLimits(constant.MinimumScreenWidth, constant.MinimumScreenHeight, -1, -1)
	ebiten.SetWindowTitle(constant.WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(true)
	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package game

import (
	"fmt"
	"image/color"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/aethiopicuschan/odori/io"
	"github.com/aethiopicuschan/odori/sprite"
	"github.com/aethiopicuschan/odori/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	components   []ui.Component
	exportButton *ui.Button
	noticer      *ui.Noticer
	explorer     *ui.Explorer
	animation    *ui.Animation
}

func NewGame() ebiten.Game {
	game := &Game{}

	buttonOffset := 10
	buttonWidth := constant.MenuWidth - buttonOffset*2
	buttonHeight := 30
	buttonMap := map[string]func(){}
	buttonMap["Load files"] = game.loadFiles
	buttonMap["Load sprite sheet"] = game.loadSpriteSheet
	buttonMap["Export"] = game.export
	buttonList := []string{
		"Load files",
		"Load sprite sheet",
		"Export",
	}
	buttons := []ui.Component{}
	i := 0
	for _, name := range buttonList {
		button := ui.NewButton(buttonOffset, buttonOffset*(i+1)+buttonHeight*i, buttonWidth, buttonHeight, name, buttonMap[name])
		if name == "Export" {
			button.SetDisabled(true)
			game.exportButton = button
		}
		buttons = append(buttons, button)
		i++
	}
	menu := ui.NewMenu(buttons)
	game.components = append(game.components, menu)

	game.explorer = ui.NewExplorer(func(s sprite.Sprite) {
		game.animation.Append(s)
	})
	game.components = append(game.components, game.explorer)

	noticeHeight := 30
	game.noticer = ui.NewNoticer(noticeHeight)

	game.animation = ui.NewAnimation(game.noticer, game.changeAnimationSize)
	game.components = append(game.components, game.animation)

	game.components = append(game.components, game.noticer)
	return game
}

func (g *Game) Update() error {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	if g.animation.CanExport() {
		g.exportButton.SetDisabled(false)
	} else {
		g.exportButton.SetDisabled(true)
	}
	for _, c := range g.components {
		c.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, c := range g.components {
		c.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nTPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	w, h := ebiten.WindowSize()
	if w != outsideWidth || h != outsideHeight {
		ebiten.SetWindowSize(outsideWidth, outsideHeight)
		for _, c := range g.components {
			c.Layout(outsideWidth, outsideHeight)
		}
	}
	return outsideWidth, outsideHeight
}

func (g *Game) loadFiles() {
	go func() {
		chPick := make(chan io.PickResult)
		go io.PickMultiple(chPick)
		result := <-chPick
		close(chPick)
		if result.Err != nil {
			if result.Err.Error() != "dialog canceled" {
				g.noticer.AddNotice(ui.WARN, result.Err.Error())
			}
			return
		}
		readCh := make(chan io.ReadResult, len(result.Paths))
		for _, path := range result.Paths {
			go io.Read(readCh, path)
		}
		appended := 0
		for _, path := range result.Paths {
			result := <-readCh
			if result.Err != nil {
				g.noticer.AddNotice(ui.ERROR, fmt.Sprintf("%s: %s", result.Err.Error(), path))
				continue
			}
			g.explorer.AppendSprite(result.Sprite)
			appended++
		}
		if appended == 0 {
			g.noticer.AddNotice(ui.WARN, "No sprite is loaded!")
		} else {
			g.noticer.AddNotice(ui.INFO, fmt.Sprintf("%d sprites are loaded!", appended))
		}
	}()
}

func (g *Game) loadSpriteSheet() {
	go func() {
		chPick := make(chan io.PickResult)
		go io.Pick(chPick)
		result := <-chPick
		close(chPick)
		if result.Err != nil {
			if result.Err.Error() != "dialog canceled" {
				g.noticer.AddNotice(ui.WARN, result.Err.Error())
			}
			return
		}
		chRead := make(chan io.ReadSpriteSheetResult)
		go io.ReadSpriteSheet(chRead, result.Paths[0])
		readResult := <-chRead
		close(chRead)
		if readResult.Err != nil {
			g.noticer.AddNotice(ui.ERROR, fmt.Sprintf("%s: %s", readResult.Err.Error(), result.Paths[0]))
			return
		}
		for _, sprite := range readResult.Sprites {
			g.explorer.AppendSprite(sprite)
		}
		appended := len(readResult.Sprites)
		if appended == 0 {
			g.noticer.AddNotice(ui.WARN, "No sprite is loaded!")
		} else {
			g.noticer.AddNotice(ui.INFO, fmt.Sprintf("%d sprites are loaded!", appended))
		}
	}()
}

func (g *Game) changeAnimationSize() {
	go func() {
		ch := make(chan io.EntryResult)
		go io.Entry(ch, "Change animation size", "Enter the size of animation in pixel", fmt.Sprintf("%dx%d", g.animation.AnimationWidth, g.animation.AnimationHeight))
		result := <-ch
		close(ch)
		if result.Err != nil {
			if result.Err.Error() != "dialog canceled" {
				g.noticer.AddNotice(ui.ERROR, result.Err.Error())
			}
			return
		}
		var animationWidth, animationHeight int
		_, err := fmt.Sscanf(result.Input, "%dx%d", &animationWidth, &animationHeight)
		if err != nil {
			g.noticer.AddNotice(ui.ERROR, err.Error())
			return
		}
		if animationWidth <= 0 || animationHeight <= 0 {
			g.noticer.AddNotice(ui.ERROR, "Invalid size!")
			return
		}
		g.animation.AnimationWidth = animationWidth
		g.animation.AnimationHeight = animationHeight
		g.noticer.AddNotice(ui.INFO, fmt.Sprintf("Animation size is changed to %dx%d", animationWidth, animationHeight))
	}()
}

func (g *Game) export() {
	g.noticer.AddNotice(ui.ERROR, "Not implemented yet!")
}

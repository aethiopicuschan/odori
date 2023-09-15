package game

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"

	"github.com/aethiopicuschan/odori/animation"
	"github.com/aethiopicuschan/odori/constant"
	"github.com/aethiopicuschan/odori/io"
	"github.com/aethiopicuschan/odori/sprite"
	"github.com/aethiopicuschan/odori/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	components []ui.Component
	buttons    []*ui.Button
	noticer    *ui.Noticer
	explorer   *ui.Explorer
	player     *ui.Player
	name       string
}

func NewGame() ebiten.Game {
	game := &Game{}

	buttonOffset := 10
	buttonWidth := constant.MenuWidth - buttonOffset*2
	buttonHeight := 30
	buttonMap := map[string]func(){}
	buttonMap["New animation"] = game.newAnimation
	buttonMap["Load files"] = game.loadFiles
	buttonMap["Load sprite sheet"] = game.loadSpriteSheet
	buttonMap["Export"] = game.export
	buttonList := []string{
		"New animation",
		"Export",
		"Load files",
		"Load sprite sheet",
	}
	buttons := []ui.Component{}
	i := 0
	for _, name := range buttonList {
		button := ui.NewButton(buttonOffset, buttonOffset*(i+1)+buttonHeight*i, buttonWidth, buttonHeight, name, buttonMap[name])
		game.buttons = append(game.buttons, button)
		if name != "New animation" {
			button.SetDisabled(true)
		}
		buttons = append(buttons, button)
		i++
	}
	menu := ui.NewMenu(buttons)
	game.components = append(game.components, menu)

	game.noticer = ui.NewNoticer()
	game.components = append(game.components, game.noticer)

	return game
}

func (g *Game) Update() error {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	for _, button := range g.buttons {
		if button.Label() == "Export" {
			if g.player != nil && g.player.CanExport() {
				button.SetDisabled(false)
			} else {
				button.SetDisabled(true)
			}
		}
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

func (g *Game) newAnimation() {
	if g.name != "" {
		return
	}
	ch := make(chan io.EntryResult)
	go io.Entry(ch, "New Animation", "Enter the project name of your new animation", "animation")
	result := <-ch
	close(ch)
	if result.Err != nil {
		if result.Err.Error() != "dialog canceled" {
			g.noticer.AddNotice(ui.ERROR, result.Err.Error())
		}
		return
	}
	if result.Input == "" {
		return
	}
	g.name = result.Input
	ebiten.SetWindowTitle(constant.WindowTitle + " - " + g.name)
	for _, button := range g.buttons {
		if button.Label() == "New animation" {
			button.SetDisabled(true)
		} else if button.Label() != "Export" {
			button.SetDisabled(false)
		}
	}
	for i, component := range g.components {
		// noticerを一度消す
		if component == g.noticer {
			g.components = append(g.components[:i], g.components[i+1:]...)
			break
		}
	}
	g.explorer = ui.NewExplorer(func(s sprite.Sprite) {
		g.player.Append(s)
	})
	g.components = append(g.components, g.explorer)
	g.player = ui.NewPlayer(g.noticer, g.changeAnimationSize)
	g.components = append(g.components, g.player)
	g.components = append(g.components, g.noticer)
}

func (g *Game) loadFiles() {
	go func() {
		pickCh := make(chan io.PickResult)
		go io.PickMultiple(pickCh)
		result := <-pickCh
		close(pickCh)
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
		for i := 0; i < cap(readCh); i++ {
			result := <-readCh
			if result.Err != nil {
				g.noticer.AddNotice(ui.ERROR, fmt.Sprintf("%s: %s", result.Err.Error(), result.Path))
				continue
			}
			g.explorer.AppendSprite(result.Sprite)
			appended++
		}
		close(readCh)
		if appended == 0 {
			g.noticer.AddNotice(ui.WARN, "No sprite is loaded!")
		} else {
			g.noticer.AddNotice(ui.INFO, fmt.Sprintf("%d sprites are loaded!", appended))
		}
	}()
}

func (g *Game) loadSpriteSheet() {
	go func() {
		pickCh := make(chan io.PickResult)
		go io.Pick(pickCh)
		result := <-pickCh
		close(pickCh)
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
		raw := g.player.RawAnimation()
		ch := make(chan io.EntryResult)
		go io.Entry(ch, "Change animation size", "Enter the size of animation in pixel", fmt.Sprintf("%dx%d", raw.Width, raw.Height))
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
		raw.Width = animationWidth
		raw.Height = animationHeight
		g.noticer.AddNotice(ui.INFO, fmt.Sprintf("Animation size is changed to %dx%d", animationWidth, animationHeight))
	}()
}

// TODO パスと名前を指定できるようにする
func (g *Game) export() {
	if !g.player.CanExport() {
		return
	}
	go func() {
		raw := g.player.RawAnimation()
		m := map[string]sprite.Sprite{}
		for _, part := range raw.Parts {
			if !part.Sprite.IsEmpty() {
				m[part.Sprite.Id()] = part.Sprite
			}
		}
		sprites := []sprite.Sprite{}
		for _, sprite := range m {
			sprites = append(sprites, sprite)
		}
		spriteSheet := map[string]image.Point{}
		// スプライトシートの出力
		if len(sprites) != 0 {
			ch := make(chan io.WriteSpriteSheetResult)
			go io.WriteSpriteSheet(ch, sprites, "./spritesheet.png")
			result := <-ch
			close(ch)
			if result.Err != nil {
				g.noticer.AddNotice(ui.ERROR, result.Err.Error())
				return
			}
			spriteSheet = result.PointsMap
		}
		// AnimationのJSON出力
		bytes, err := json.MarshalIndent(animation.AnimationP{
			Name:        g.name,
			Animation:   raw,
			SpriteSheet: spriteSheet,
		}, "", "  ")
		if err != nil {
			g.noticer.AddNotice(ui.ERROR, err.Error())
			return
		}
		ch := make(chan error)
		go io.Write(ch, bytes, "./animation.json")
		err = <-ch
		close(ch)
		if err != nil {
			g.noticer.AddNotice(ui.ERROR, err.Error())
			return
		}
		g.noticer.AddNotice(ui.INFO, "Exported!")
	}()
}

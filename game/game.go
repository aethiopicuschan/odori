package game

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"path"
	"path/filepath"

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
	buttonMap["Import"] = game.importAnimation
	buttonMap["Export"] = game.exportAnimation
	buttonMap["Export as GIF"] = game.exportAsGif
	buttonList := []string{
		"New animation",
		"Import",
		"Export",
		"Export as GIF",
		"Load files",
		"Load sprite sheet",
	}
	buttons := []ui.Component{}
	i := 0
	for _, name := range buttonList {
		button := ui.NewButton(buttonOffset, buttonOffset*(i+1)+buttonHeight*i, buttonWidth, buttonHeight, name, buttonMap[name])
		game.buttons = append(game.buttons, button)
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
	// TODO 開いているプロジェクトがあるときにWindowを閉じるときは確認する
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	for _, button := range g.buttons {
		if g.name == "" {
			button.SetDisabled(button.Label() != "New animation" && button.Label() != "Import")
		} else {
			if button.Label() == "Export" || button.Label() == "Export as GIF" {
				button.SetDisabled(!g.player.RawAnimation().CanExport())
			} else {
				button.SetDisabled(button.Label() == "New animation" || button.Label() == "Import")
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

// newAnimationとImportから呼ばれる想定
func (g *Game) startProject(name string) {
	if !animation.IsValidName(name) {
		g.noticer.AddNotice(ui.ERROR, "Invalid name!")
		return
	}
	g.name = name
	ebiten.SetWindowTitle(constant.WindowTitle + " - " + g.name)
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
	funcMap := map[string]func(){
		"changeAnimationSize": g.changeAnimationSize,
		"renameAnimation":     g.renameAnimation,
	}
	g.player = ui.NewPlayer(g.name, g.noticer, funcMap)
	g.components = append(g.components, g.player)
	g.components = append(g.components, g.noticer)
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
	g.startProject(result.Input)
}

func (g *Game) loadFiles() {
	go func() {
		pickCh := make(chan io.PickMultipleResult)
		go io.PickMultiple(pickCh, io.WithName("Select images"), io.WithPatterns([]string{"*.png"}))
		result := <-pickCh
		close(pickCh)
		if result.Err != nil {
			if result.Err.Error() != "dialog canceled" {
				g.noticer.AddNotice(ui.WARN, result.Err.Error())
			}
			return
		}
		readCh := make(chan io.ReadSpriteResult, len(result.Paths))
		for _, path := range result.Paths {
			go io.ReadSprite(readCh, path)
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
		go io.Pick(pickCh, io.WithName("Select sprite sheet"), io.WithPatterns([]string{"*.png"}))
		result := <-pickCh
		close(pickCh)
		if result.Err != nil {
			if result.Err.Error() != "dialog canceled" {
				g.noticer.AddNotice(ui.WARN, result.Err.Error())
			}
			return
		}
		chRead := make(chan io.ReadSpriteSheetResult)
		go io.ReadSpriteSheet(chRead, result.Path)
		readResult := <-chRead
		close(chRead)
		if readResult.Err != nil {
			g.noticer.AddNotice(ui.ERROR, fmt.Sprintf("%s: %s", readResult.Err.Error(), result.Path))
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

func (g *Game) renameAnimation() {
	go func() {
		ch := make(chan io.EntryResult)
		go io.Entry(ch, "Change animation name", "Enter new name", g.name)
		result := <-ch
		close(ch)
		if result.Err != nil {
			if result.Err.Error() != "dialog canceled" {
				g.noticer.AddNotice(ui.ERROR, result.Err.Error())
			}
			return
		}
		if !animation.IsValidName(result.Input) {
			g.noticer.AddNotice(ui.ERROR, "Invalid name!")
			return
		}
		g.name = result.Input
		g.player.Rename(g.name)
		ebiten.SetWindowTitle(constant.WindowTitle + " - " + g.name)
		g.noticer.AddNotice(ui.INFO, fmt.Sprintf(`Animation name is changed to "%s"`, g.name))
	}()
}

func (g *Game) exportAnimation() {
	if !g.player.RawAnimation().CanExport() {
		return
	}
	g.player.Stop()
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
	selectDirCh := make(chan io.SelectDirResult)
	go io.SelectDir(selectDirCh)
	result := <-selectDirCh
	close(selectDirCh)
	if result.Err != nil {
		if result.Err.Error() != "dialog canceled" {
			g.noticer.AddNotice(ui.ERROR, result.Err.Error())
		}
		return
	}
	dir := result.Path
	spriteSheetPath := filepath.Join(dir, g.name+".png")
	jsonPath := filepath.Join(dir, g.name+".json")
	if io.IsExist(spriteSheetPath) || io.IsExist(jsonPath) {
		questionCh := make(chan io.QuestionResult)
		go io.Question(questionCh, "Overwrite", "Overwrite existing files?")
		result := <-questionCh
		close(questionCh)
		if !result.Answer {
			return
		}
	}
	spriteSheet := map[string]image.Rectangle{}
	// スプライトシートの出力
	if len(sprites) != 0 {
		ch := make(chan io.WriteSpriteSheetResult)
		go io.WriteSpriteSheet(ch, sprites, filepath.Join(dir, g.name+".png"))
		result := <-ch
		close(ch)
		if result.Err != nil {
			g.noticer.AddNotice(ui.ERROR, result.Err.Error())
			return
		}
		spriteSheet = result.RectsMap
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
	writeCh := make(chan error)
	go io.Write(writeCh, bytes, filepath.Join(dir, g.name+".json"))
	err = <-writeCh
	close(writeCh)
	if err != nil {
		g.noticer.AddNotice(ui.ERROR, err.Error())
		return
	}
	g.noticer.AddNotice(ui.INFO, "Exported!")
}

func (g *Game) importAnimation() {
	// JSONを読み込ませる
	pickCh := make(chan io.PickResult)
	go io.Pick(pickCh, io.WithName("Select animation"), io.WithPatterns([]string{"*.json"}))
	result := <-pickCh
	close(pickCh)
	if result.Err != nil {
		if result.Err.Error() != "dialog canceled" {
			g.noticer.AddNotice(ui.WARN, result.Err.Error())
		}
		return
	}
	// JSONの読み込み
	readCh := make(chan io.ReadResult)
	go io.Read(readCh, result.Path)
	readResult := <-readCh
	close(readCh)
	if readResult.Err != nil {
		g.noticer.AddNotice(ui.ERROR, fmt.Sprintf("%s: %s", readResult.Err.Error(), result.Path))
		return
	}
	var animationP animation.AnimationP
	err := json.Unmarshal(readResult.Bytes, &animationP)
	if err != nil {
		g.noticer.AddNotice(ui.ERROR, err.Error())
		return
	}
	withSpriteSheet := false
	for _, part := range animationP.Animation.Parts {
		if !part.Sprite.IsEmpty() {
			withSpriteSheet = true
			break
		}
	}
	// スプライトシートの読み込み
	if withSpriteSheet {
		si, err := io.ReadPng(path.Join(path.Dir(result.Path), animationP.Name+".png"))
		if err != nil {
			g.noticer.AddNotice(ui.ERROR, err.Error())
			return
		}
		sprites := sprite.NewSpritesFromRectMap(si, animationP.SpriteSheet)
		g.startProject(animationP.Name)
		for _, sprite := range sprites {
			g.explorer.AppendSprite(sprite)
		}
		for i, part := range animationP.Animation.Parts {
			if !part.Sprite.IsEmpty() {
				for _, sprite := range sprites {
					if sprite.Id() == part.Sprite.Id() {
						animationP.Animation.Parts[i].Sprite = sprite
						break
					}
				}
			}
		}
		g.player.Import(animationP.Animation)
		g.noticer.AddNotice(ui.INFO, fmt.Sprintf(`Project "%s" was imported with %d sprites!`, animationP.Name, len(sprites)))
	} else {
		g.startProject(animationP.Name)
		g.player.Import(animationP.Animation)
		g.noticer.AddNotice(ui.INFO, fmt.Sprintf(`Project "%s" was imported!`, animationP.Name))
	}
}

func (g *Game) exportAsGif() {
	if !g.player.RawAnimation().CanExport() {
		return
	}
	pickCh := make(chan io.PickResult)
	go io.Pick(pickCh, io.WithName("Save as..."), io.WithPatterns([]string{"*.gif"}), io.WithToSave(g.name+".gif"))
	result := <-pickCh
	close(pickCh)
	if result.Err != nil {
		if result.Err.Error() != "dialog canceled" {
			g.noticer.AddNotice(ui.WARN, result.Err.Error())
		}
		return
	}
	err := g.player.RawAnimation().ExportAsGif(result.Path)
	if err != nil {
		g.noticer.AddNotice(ui.ERROR, err.Error())
		return
	}
	g.noticer.AddNotice(ui.INFO, "Exported!")
}

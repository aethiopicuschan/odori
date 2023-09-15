package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/aethiopicuschan/odori/animation"
	"github.com/aethiopicuschan/odori/constant"
	"github.com/aethiopicuschan/odori/io"
	"github.com/aethiopicuschan/odori/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type Animation struct {
	width               int
	height              int
	offsetY             int
	offsetX             int
	animation           *animation.Animation
	font                font.Face
	customLinks         []*Link
	parts               []*animation.Part
	indexes             []int
	currentPart         int
	playing             bool
	currentTick         int
	maxTick             int
	playerButtons       []*Button
	playerX             int
	playerY             int
	playerWidth         int
	playerHeight        int
	cursolOnPlayer      bool
	noticer             *Noticer
	changeAnimationSize func()
}

func NewAnimation(noticer *Noticer, changeAnimationSize func()) *Animation {
	w, h := ebiten.WindowSize()
	tt, _ := opentype.Parse(goregular.TTF)
	font, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})

	a := &Animation{}
	a.noticer = noticer
	a.animation = animation.NewAnimation()
	a.changeAnimationSize = changeAnimationSize
	a.offsetX = constant.MenuWidth
	a.offsetY = h / 3
	a.width = w - a.offsetX
	a.height = h - a.offsetY
	a.font = font
	a.parts = []*animation.Part{}
	a.currentPart = -1
	a.indexes = []int{}
	a.playing = false
	a.currentTick = 0
	a.maxTick = 0
	a.playerX = 0
	a.playerY = 0
	a.playerWidth = 0
	a.playerHeight = 0
	a.cursolOnPlayer = false
	a.playerButtons = []*Button{
		NewButton(0, 0, 20, 30, "<<", func() {
			if a.playing {
				return
			}
			if a.currentPart == 0 {
				a.currentTick = 0
			} else if a.currentPart > 0 {
				a.currentTick = a.indexes[a.currentPart-1]
				a.currentPart -= 1
			}
		}),
		NewButton(0, 0, 20, 30, "<", func() {
			if a.playing {
				return
			}
			a.currentTick--
			if a.currentTick < 0 {
				a.currentTick = a.maxTick - 1
			}
		}),
		NewButton(0, 0, 40, 30, "Play", func() {
			if a.maxTick == 0 {
				a.playing = false
				return
			}
			a.playing = !a.playing
		}),
		NewButton(0, 0, 20, 30, ">", func() {
			if a.playing {
				return
			}
			a.currentTick++
			if a.currentTick >= a.maxTick {
				a.currentTick = 0
			}
		}),
		NewButton(0, 0, 20, 30, ">>", func() {
			if a.playing {
				return
			}
			if a.currentPart+1 < len(a.parts) {
				a.currentTick = a.indexes[a.currentPart+1]
				a.currentPart += 1
			}
		}),
	}
	a.customLinks = []*Link{
		NewLink(0, 0, "# Informations", nil),
		NewLink(0, 0, "Size", a.changeAnimationSize),
		NewLink(0, 0, "TPS", nil),
		NewLink(0, 0, "TotalLen", nil),
		NewLink(0, 0, "# Properties", nil),
		NewLink(0, 0, "Scale", func() {
			a.playing = false
			go func() {
				if len(a.parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := a.parts[a.currentPart]
				go io.Entry(ch, "Change scale", "Enter the scale", fmt.Sprintf("%0.2f", part.Scale))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						a.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var scale float64
				_, err := fmt.Sscanf(result.Input, "%f", &scale)
				if err != nil {
					a.noticer.AddNotice(ERROR, err.Error())
					return
				}
				if scale <= 0 {
					a.noticer.AddNotice(ERROR, "Scale must be greater than 0")
					return
				}
				part.Scale = scale
			}()
		}),
		NewLink(0, 0, "DiffX", func() {
			a.playing = false
			go func() {
				if len(a.parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := a.parts[a.currentPart]
				go io.Entry(ch, "Change DiffX", "Enter the diff of position from center.", fmt.Sprintf("%d", part.DiffX))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						a.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var diffX int
				_, err := fmt.Sscanf(result.Input, "%d", &diffX)
				if err != nil {
					a.noticer.AddNotice(ERROR, err.Error())
					return
				}
				part.DiffX = diffX
			}()
		}),
		NewLink(0, 0, "DiffY", func() {
			a.playing = false
			go func() {
				if len(a.parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := a.parts[a.currentPart]
				go io.Entry(ch, "Change DiffY", "Enter the diff of position from center.", fmt.Sprintf("%d", part.DiffY))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						a.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var diffY int
				_, err := fmt.Sscanf(result.Input, "%d", &diffY)
				if err != nil {
					a.noticer.AddNotice(ERROR, err.Error())
					return
				}
				part.DiffY = diffY
			}()
		}),
		NewLink(0, 0, "Reverse", func() {
			a.playing = false
			if len(a.parts) == 0 {
				return
			}
			part := a.parts[a.currentPart]
			part.Reverse = !part.Reverse
		}),
		NewLink(0, 0, "Len", func() {
			a.playing = false
			go func() {
				if len(a.parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := a.parts[a.currentPart]
				go io.Entry(ch, "Change Length", fmt.Sprintf("Enter the length of ticks. (%d ticks means 1sec)", ebiten.TPS()), fmt.Sprintf("%d", part.Length))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						a.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var length int
				_, err := fmt.Sscanf(result.Input, "%d", &length)
				if err != nil {
					a.noticer.AddNotice(ERROR, err.Error())
					return
				}
				if length <= 0 {
					a.noticer.AddNotice(ERROR, "Length must be greater than 0")
					return
				}
				part.Length = length
				// lengthが変わった関係で諸々のパラメータをリセットする必要がある
				a.resetIndexes()
			}()
		}),
		NewLink(0, 0, "# Operations", nil),
		NewLink(0, 0, "Auto scale", func() {
			a.playing = false
			if len(a.parts) == 0 {
				return
			}
			part := a.parts[a.currentPart]
			if part.Sprite.IsEmpty() {
				return
			}
			scale := 1.0
			width := part.Sprite.Image.Bounds().Dx()
			height := part.Sprite.Image.Bounds().Dy()
			if width < a.animation.Width && height < a.animation.Height {
				if width > height {
					scale = float64(a.animation.Width / width)
				} else {
					scale = float64(a.animation.Height / height)
				}
			}
			if width > a.animation.Width || height > a.animation.Height {
				if width > height {
					scale = float64(a.animation.Width) / float64(width)
				} else {
					scale = float64(a.animation.Height) / float64(height)
				}
			}
			part.Scale = scale
		}),
		NewLink(0, 0, "Reset", func() {
			a.playing = false
			go func() {
				if len(a.parts) == 0 {
					return
				}
				part := a.parts[a.currentPart]
				ch := make(chan io.QuestionResult)
				go io.Question(ch, "Reset", "Are you sure you want to reset properties of current part?")
				result := <-ch
				close(ch)
				if !result.Answer {
					return
				}
				part.Scale = 1.0
				part.DiffX = 0
				part.DiffY = 0
				part.Reverse = false
				part.Length = ebiten.TPS()
				a.resetIndexes()
			}()
		}),
		NewLink(0, 0, "Delete", func() {
			a.playing = false
			go func() {
				if len(a.parts) == 0 {
					return
				}
				ch := make(chan io.QuestionResult)
				go io.Question(ch, "Delete", "Are you sure you want to delete current part?")
				result := <-ch
				close(ch)
				if !result.Answer {
					return
				}
				if len(a.parts) == 1 {
					a.parts = []*animation.Part{}
					a.currentPart = -1
				} else {
					a.parts = append(a.parts[:a.currentPart], a.parts[a.currentPart+1:]...)
					a.currentPart--
					if a.currentPart < 0 {
						a.currentPart = 0
					}
				}
				a.resetIndexes()
			}()
		}),
	}

	a.Layout(w, h)

	return a
}

func (a *Animation) resetIndexes() {
	a.maxTick = 0
	a.currentTick = 0
	a.indexes = []int{}
	for i, p := range a.parts {
		if i == a.currentPart {
			a.currentTick = a.maxTick
		}
		a.indexes = append(a.indexes, a.maxTick)
		a.maxTick += p.Length
	}
}

func (a *Animation) Append(sprite sprite.Sprite) {
	if len(a.parts) == 0 {
		a.currentPart = 0
	}
	len := ebiten.TPS()
	a.parts = append(a.parts, animation.NewPart(sprite, len))
	a.indexes = append(a.indexes, a.maxTick)
	a.currentTick = a.maxTick
	a.maxTick += len
}

func (a *Animation) Update() error {
	partsLen := len(a.parts)
	for _, b := range a.playerButtons {
		b.SetDisabled(partsLen == 0 || a.playing)
		if b.label == "Play" || b.label == "Stop" {
			b.SetDisabled(partsLen == 0)
			if a.playing {
				b.label = "Stop"
			} else {
				b.label = "Play"
			}
		}
		b.Update()
	}
	for _, l := range a.customLinks {
		if l.id == "Size" {
			l.SetLabel(fmt.Sprintf("Size: %dx%d", a.animation.Width, a.animation.Height))
		}
		if l.id == "TPS" {
			l.SetLabel(fmt.Sprintf("TPS : %d", ebiten.TPS()))
		}
		if l.id == "TotalLen" {
			l.SetLabel(fmt.Sprintf("Len : %d ticks (%0.2f sec)", a.maxTick, float64(a.maxTick)/float64(ebiten.TPS())))
		}
		if l.id == "# Properties" {
			break
		} else {
			if partsLen == 0 {
				l.Update()
			}
		}
	}

	// Tickと索引を元に現在のパーツを更新
	if partsLen > 0 {
		if a.currentTick == 0 {
			a.currentPart = 0
		}
		if a.currentPart+1 < partsLen && a.currentTick >= a.indexes[a.currentPart+1] {
			a.currentPart++
		} else if a.currentPart > 0 && a.currentTick < a.indexes[a.currentPart] {
			a.currentPart--
		}
		part := a.parts[a.currentPart]
		// カスタムリンク
		for _, l := range a.customLinks {
			if l.id == "Scale" {
				l.SetLabel(fmt.Sprintf("Scale: %0.2f", part.Scale))
			}
			if l.id == "DiffX" {
				l.SetLabel(fmt.Sprintf("DiffX: %d", part.DiffX))
			}
			if l.id == "DiffY" {
				l.SetLabel(fmt.Sprintf("DiffY: %d", part.DiffY))
			}
			if l.id == "Reverse" {
				l.SetLabel(fmt.Sprintf("Reverse: %t", part.Reverse))
			}
			if l.id == "Len" {
				l.SetLabel(fmt.Sprintf("Len  : %d ticks (%0.2f sec)", part.Length, float64(part.Length)/float64(ebiten.TPS())))
			}
			l.SetDisabled(a.playing)
			if part.Sprite.IsEmpty() && l.id != "Len" && l.id != "Reset" && l.id != "Delete" && l.id != "Size" {
				l.SetDisabled(true)
			}
			l.Update()
		}
	}
	// 再生中ならTickを進める
	if a.playing {
		a.currentTick++
		if a.currentTick >= a.maxTick {
			a.currentTick = 0
		}
	}

	a.cursolOnPlayer = false
	if !ebiten.IsFocused() {
		return nil
	}

	if len(a.parts) == 0 {
		return nil
	}

	// Play/Stop by space key.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		a.playing = !a.playing
	}

	// Arrow keys
	if !a.playing {
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			a.currentTick--
			if a.currentTick < 0 {
				a.currentTick = 0
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			a.currentTick++
			if a.currentTick >= a.maxTick {
				a.currentTick = a.maxTick - 1
			}
		}
	}

	// Cursor
	cursorX, cursorY := ebiten.CursorPosition()
	isCursorOnAnimation := cursorX >= a.offsetX && cursorX <= a.offsetX+a.width && cursorY >= a.offsetY && cursorY <= a.offsetY+a.height
	if !isCursorOnAnimation {
		return nil
	}
	xOnAnimation := cursorX - a.offsetX
	yOnAnimation := cursorY - a.offsetY
	isCurorOnPlayer := xOnAnimation >= a.playerX && xOnAnimation <= a.playerX+a.playerWidth && yOnAnimation >= a.playerY && yOnAnimation <= a.playerY+a.playerHeight
	if !isCurorOnPlayer {
		return nil
	}
	a.cursolOnPlayer = true
	ebiten.SetCursorShape(ebiten.CursorShapePointer)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		a.playing = false
		a.currentTick = int(float64(xOnAnimation-a.playerX) / float64(a.playerWidth) * float64(a.maxTick))
		if a.currentTick >= a.maxTick {
			a.currentTick = a.maxTick - 1
		}
		for i, index := range a.indexes {
			if a.currentTick >= index {
				a.currentPart = i
			} else {
				break
			}
		}
	}

	return nil
}

func (a *Animation) Draw(screen *ebiten.Image) {
	// Fill background.
	bg := ebiten.NewImage(screen.Bounds().Dx()-constant.MenuWidth, a.height)
	bgOp := &ebiten.DrawImageOptions{}
	bgOp.GeoM.Translate(float64(constant.MenuWidth), float64(a.offsetY))

	// Play/Stop button.
	for _, b := range a.playerButtons {
		defer b.Draw(screen)
	}
	for _, cl := range a.customLinks {
		if cl.id == "# Properties" && len(a.parts) == 0 {
			break
		}
		defer cl.Draw(screen)
	}

	// Final Draw.
	defer screen.DrawImage(bg, bgOp)

	// Animation Frame.
	frameLine := ebiten.NewImage(a.animation.Width+2, a.animation.Height+2)
	frameLine.Fill(color.Black)
	frame := ebiten.NewImage(a.animation.Width, a.animation.Height)
	frame.Fill(color.White)
	{
		isGray := false
		boxSize := constant.DefaultAnimationSize / 5
		maxX := int(math.Ceil(float64(a.animation.Width) / float64(boxSize)))
		maxY := int(math.Ceil(float64(a.animation.Height) / float64(boxSize)))
		boxGray := ebiten.NewImage(boxSize, boxSize)
		boxGray.Fill(color.Gray{Y: constant.ExplorerGrayY})
		boxTransparent := ebiten.NewImage(boxSize, boxSize)
		boxTransparent.Fill(color.Transparent)

		for x := 0; x < maxX; x++ {
			for y := 0; y < maxY; y++ {
				boxOp := &ebiten.DrawImageOptions{}
				boxOp.GeoM.Translate(float64(x*boxSize), float64(y*boxSize))
				if isGray {
					frame.DrawImage(boxGray, boxOp)
				} else {
					frame.DrawImage(boxTransparent, boxOp)
				}
				isGray = !isGray
			}
			if maxY%2 == 0 {
				isGray = !isGray
			}
		}
	}
	x := float64(bg.Bounds().Dx()/2) - float64(a.animation.Width)/2
	y := float64(bg.Bounds().Dy()/3) - float64(a.animation.Height)/2
	frameOp := &ebiten.DrawImageOptions{}
	frameLineOp := &ebiten.DrawImageOptions{}
	frameOp.GeoM.Translate(float64(x), float64(y))
	frameLineOp.GeoM.Translate(float64(x-1), float64(y-1))
	bg.DrawImage(frameLine, frameLineOp)

	// Draw part.
	if a.currentPart >= 0 {
		part := a.parts[a.currentPart]
		scale := part.Scale
		diffX := part.DiffX
		diffY := part.DiffY
		reverse := part.Reverse
		if !part.Sprite.IsEmpty() {
			img := part.Sprite.Image
			op := &ebiten.DrawImageOptions{}
			if reverse {
				op.GeoM.Scale(-1, 1)
				op.GeoM.Translate(float64(img.Bounds().Dx()), 0)
				diffX *= -1
			}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(-((float64(img.Bounds().Dx()-diffX))*scale)/2, -((float64(img.Bounds().Dy()-diffY))*scale)/2)
			op.GeoM.Translate(float64(a.animation.Width/2), float64(a.animation.Height/2))
			frame.DrawImage(img, op)
		}
	}

	bg.DrawImage(frame, frameOp)

	// Player.
	player := ebiten.NewImage(a.playerWidth, a.playerHeight)
	player.Fill(color.Gray{Y: constant.PlayerGrayY})
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(10, float64(a.playerY))
	prs := int(float64(a.currentTick) / float64(a.maxTick) * float64(player.Bounds().Dx()))
	progress := ebiten.NewImage(2, player.Bounds().Dy())
	progress.Fill(color.Black)
	progressOp := &ebiten.DrawImageOptions{}
	progressOp.GeoM.Translate(float64(prs), 0)
	// draw index positions
	indexPos := ebiten.NewImage(2, player.Bounds().Dy())
	indexPos.Fill(color.Gray{Y: constant.PlayerIndexGrayY})
	for i, index := range a.indexes {
		x := int(float64(index) / float64(a.maxTick) * float64(player.Bounds().Dx()))
		indexPosOp := &ebiten.DrawImageOptions{}
		indexPosOp.GeoM.Translate(float64(x), 0)
		player.DrawImage(indexPos, indexPosOp)
		if i == a.currentPart {
			// Fill current part.
			width := int(math.Ceil(float64(a.parts[i].Length) / float64(a.maxTick) * float64(player.Bounds().Dx())))
			currentPart := ebiten.NewImage(width, player.Bounds().Dy())
			currentPart.Fill(color.Gray{Y: constant.PlayerCurrentGrayY})
			currentPartOp := &ebiten.DrawImageOptions{}
			currentPartOp.GeoM.Translate(float64(x), 0)
			player.DrawImage(currentPart, currentPartOp)
		}
	}
	player.DrawImage(progress, progressOp)
	bg.DrawImage(player, playerOp)
	s := fmt.Sprintf("Part: %d / %d\nTick: %d / %d\nSec: %0.2f / %0.2f", a.currentPart+1, len(a.parts), a.currentTick, a.maxTick, float64(a.currentTick)/float64(ebiten.TPS()), float64(a.maxTick)/float64(ebiten.TPS()))
	bs := text.BoundString(a.font, s)
	text.Draw(bg, s, a.font, 10, a.playerY-bs.Dy(), color.Black)
}

func (a *Animation) Layout(outsideWidth, outsideHeight int) {
	// Resize
	a.height = outsideHeight - outsideHeight/3
	a.offsetY = outsideHeight / 3
	a.offsetX = constant.MenuWidth
	a.width = outsideWidth - a.offsetX
	a.playerWidth = a.width - 20
	a.playerHeight = 30
	a.playerX = 10
	a.playerY = a.height - a.playerHeight - 10

	totalWidth := 0
	for _, b := range a.playerButtons {
		totalWidth += b.width + 5
	}
	startX := a.offsetX + (a.width-totalWidth)/2 + 5
	for _, b := range a.playerButtons {
		b.MoveTo(startX, a.offsetY+a.height-75)
		startX += b.width + 5
	}
	defaultBs := text.BoundString(a.font, "DEFAULT")
	for i, l := range a.customLinks {
		l.MoveTo(a.offsetX+5, a.offsetY+5+((defaultBs.Dy()+5)*i))
	}
}

func (a *Animation) CanExport() bool {
	return len(a.parts) > 0
}

func (a *Animation) RawAnimation() *animation.Animation {
	return a.animation
}

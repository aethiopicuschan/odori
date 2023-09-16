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

type Player struct {
	width               int
	height              int
	offsetY             int
	offsetX             int
	animation           *animation.Animation
	font                font.Face
	customLinks         []*Link
	indexes             []int
	currentPart         int
	playing             bool
	currentTick         int
	maxTick             int
	buttons             []*Button
	barX                int
	barY                int
	barWidth            int
	barHeight           int
	cursolOnBar         bool
	noticer             *Noticer
	changeAnimationSize func()
}

func NewPlayer(noticer *Noticer, changeAnimationSize func()) *Player {
	w, h := ebiten.WindowSize()
	tt, _ := opentype.Parse(goregular.TTF)
	font, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})

	p := &Player{}
	p.noticer = noticer
	p.animation = animation.NewAnimation()
	p.changeAnimationSize = changeAnimationSize
	p.offsetX = constant.MenuWidth
	p.offsetY = h / 3
	p.width = w - p.offsetX
	p.height = h - p.offsetY
	p.font = font
	p.currentPart = -1
	p.indexes = []int{}
	p.playing = false
	p.currentTick = 0
	p.maxTick = 0
	p.barX = 0
	p.barY = 0
	p.barWidth = 0
	p.barHeight = 0
	p.cursolOnBar = false
	p.buttons = []*Button{
		NewButton(0, 0, 20, 30, "<<", func() {
			if p.playing {
				return
			}
			if p.currentPart == 0 {
				p.currentTick = 0
			} else if p.currentPart > 0 {
				p.currentTick = p.indexes[p.currentPart-1]
				p.currentPart -= 1
			}
		}),
		NewButton(0, 0, 20, 30, "<", func() {
			if p.playing {
				return
			}
			p.currentTick--
			if p.currentTick < 0 {
				p.currentTick = p.maxTick - 1
			}
		}),
		NewButton(0, 0, 40, 30, "Play", func() {
			if p.maxTick == 0 {
				p.playing = false
				return
			}
			p.playing = !p.playing
		}),
		NewButton(0, 0, 20, 30, ">", func() {
			if p.playing {
				return
			}
			p.currentTick++
			if p.currentTick >= p.maxTick {
				p.currentTick = 0
			}
		}),
		NewButton(0, 0, 20, 30, ">>", func() {
			if p.playing {
				return
			}
			if p.currentPart+1 < len(p.animation.Parts) {
				p.currentTick = p.indexes[p.currentPart+1]
				p.currentPart += 1
			}
		}),
	}
	p.customLinks = []*Link{
		NewLink(0, 0, "# Informations", nil),
		NewLink(0, 0, "Size", p.changeAnimationSize),
		NewLink(0, 0, "TPS", nil),
		NewLink(0, 0, "TotalLen", nil),
		NewLink(0, 0, "# Properties", nil),
		NewLink(0, 0, "Scale", func() {
			p.playing = false
			go func() {
				if len(p.animation.Parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := p.animation.Parts[p.currentPart]
				go io.Entry(ch, "Change scale", "Enter the scale", fmt.Sprintf("%0.2f", part.Scale))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						p.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var scale float64
				_, err := fmt.Sscanf(result.Input, "%f", &scale)
				if err != nil {
					p.noticer.AddNotice(ERROR, err.Error())
					return
				}
				if scale <= 0 {
					p.noticer.AddNotice(ERROR, "Scale must be greater than 0")
					return
				}
				part.Scale = scale
			}()
		}),
		NewLink(0, 0, "DiffX", func() {
			p.playing = false
			go func() {
				if len(p.animation.Parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := p.animation.Parts[p.currentPart]
				go io.Entry(ch, "Change DiffX", "Enter the diff of position from center.", fmt.Sprintf("%d", part.DiffX))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						p.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var diffX int
				_, err := fmt.Sscanf(result.Input, "%d", &diffX)
				if err != nil {
					p.noticer.AddNotice(ERROR, err.Error())
					return
				}
				part.DiffX = diffX
			}()
		}),
		NewLink(0, 0, "DiffY", func() {
			p.playing = false
			go func() {
				if len(p.animation.Parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := p.animation.Parts[p.currentPart]
				go io.Entry(ch, "Change DiffY", "Enter the diff of position from center.", fmt.Sprintf("%d", part.DiffY))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						p.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var diffY int
				_, err := fmt.Sscanf(result.Input, "%d", &diffY)
				if err != nil {
					p.noticer.AddNotice(ERROR, err.Error())
					return
				}
				part.DiffY = diffY
			}()
		}),
		NewLink(0, 0, "Reverse", func() {
			p.playing = false
			if len(p.animation.Parts) == 0 {
				return
			}
			part := p.animation.Parts[p.currentPart]
			part.Reverse = !part.Reverse
		}),
		NewLink(0, 0, "Len", func() {
			p.playing = false
			go func() {
				if len(p.animation.Parts) == 0 {
					return
				}
				ch := make(chan io.EntryResult)
				part := p.animation.Parts[p.currentPart]
				go io.Entry(ch, "Change Length", fmt.Sprintf("Enter the length of ticks. (%d ticks means 1sec)", ebiten.TPS()), fmt.Sprintf("%d", part.Length))
				result := <-ch
				close(ch)
				if result.Err != nil {
					if result.Err.Error() != "dialog canceled" {
						p.noticer.AddNotice(ERROR, result.Err.Error())
					}
					return
				}
				var length int
				_, err := fmt.Sscanf(result.Input, "%d", &length)
				if err != nil {
					p.noticer.AddNotice(ERROR, err.Error())
					return
				}
				if length <= 0 {
					p.noticer.AddNotice(ERROR, "Length must be greater than 0")
					return
				}
				part.Length = length
				// lengthが変わった関係で諸々のパラメータをリセットする必要がある
				p.resetIndexes()
			}()
		}),
		NewLink(0, 0, "# Operations", nil),
		NewLink(0, 0, "Auto scale", func() {
			p.playing = false
			if len(p.animation.Parts) == 0 {
				return
			}
			part := p.animation.Parts[p.currentPart]
			if part.Sprite.IsEmpty() {
				return
			}
			scale := 1.0
			width := part.Sprite.Image.Bounds().Dx()
			height := part.Sprite.Image.Bounds().Dy()
			if width < p.animation.Width && height < p.animation.Height {
				if width > height {
					scale = float64(p.animation.Width / width)
				} else {
					scale = float64(p.animation.Height / height)
				}
			}
			if width > p.animation.Width || height > p.animation.Height {
				if width > height {
					scale = float64(p.animation.Width) / float64(width)
				} else {
					scale = float64(p.animation.Height) / float64(height)
				}
			}
			part.Scale = scale
		}),
		NewLink(0, 0, "Reset", func() {
			p.playing = false
			go func() {
				if len(p.animation.Parts) == 0 {
					return
				}
				part := p.animation.Parts[p.currentPart]
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
				p.resetIndexes()
			}()
		}),
		NewLink(0, 0, "Delete", func() {
			p.playing = false
			go func() {
				if len(p.animation.Parts) == 0 {
					return
				}
				ch := make(chan io.QuestionResult)
				go io.Question(ch, "Delete", "Are you sure you want to delete current part?")
				result := <-ch
				close(ch)
				if !result.Answer {
					return
				}
				if len(p.animation.Parts) == 1 {
					p.animation.Parts = []*animation.Part{}
					p.currentPart = -1
				} else {
					p.animation.Parts = append(p.animation.Parts[:p.currentPart], p.animation.Parts[p.currentPart+1:]...)
					p.currentPart--
					if p.currentPart < 0 {
						p.currentPart = 0
					}
				}
				p.resetIndexes()
			}()
		}),
	}

	p.Layout(w, h)

	return p
}

func (p *Player) resetIndexes() {
	p.maxTick = 0
	p.currentTick = 0
	p.indexes = []int{}
	for i, part := range p.animation.Parts {
		if i == p.currentPart {
			p.currentTick = p.maxTick
		}
		p.indexes = append(p.indexes, p.maxTick)
		p.maxTick += part.Length
	}
}

func (p *Player) Append(sprite sprite.Sprite) {
	if len(p.animation.Parts) == 0 {
		p.currentPart = 0
	}
	len := ebiten.TPS()
	p.animation.Parts = append(p.animation.Parts, animation.NewPart(sprite, len))
	p.indexes = append(p.indexes, p.maxTick)
	p.currentTick = p.maxTick
	p.maxTick += len
}

func (p *Player) Update() error {
	partsLen := len(p.animation.Parts)
	for _, b := range p.buttons {
		b.SetDisabled(partsLen == 0 || p.playing)
		if b.label == "Play" || b.label == "Stop" {
			b.SetDisabled(partsLen == 0)
			if p.playing {
				b.label = "Stop"
			} else {
				b.label = "Play"
			}
		}
		b.Update()
	}
	for _, l := range p.customLinks {
		if l.id == "Size" {
			l.SetLabel(fmt.Sprintf("Size: %dx%d", p.animation.Width, p.animation.Height))
		}
		if l.id == "TPS" {
			l.SetLabel(fmt.Sprintf("TPS : %d", ebiten.TPS()))
		}
		if l.id == "TotalLen" {
			l.SetLabel(fmt.Sprintf("Len : %d ticks (%0.2f sec)", p.maxTick, float64(p.maxTick)/float64(ebiten.TPS())))
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
		if p.currentTick == 0 {
			p.currentPart = 0
		}
		if p.currentPart+1 < partsLen && p.currentTick >= p.indexes[p.currentPart+1] {
			p.currentPart++
		} else if p.currentPart > 0 && p.currentTick < p.indexes[p.currentPart] {
			p.currentPart--
		}
		part := p.animation.Parts[p.currentPart]
		// カスタムリンク
		for _, l := range p.customLinks {
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
			l.SetDisabled(p.playing)
			if part.Sprite.IsEmpty() && l.id != "Len" && l.id != "Reset" && l.id != "Delete" && l.id != "Size" {
				l.SetDisabled(true)
			}
			l.Update()
		}
	}
	// 再生中ならTickを進める
	if p.playing {
		p.currentTick++
		if p.currentTick >= p.maxTick {
			p.currentTick = 0
		}
	}

	p.cursolOnBar = false
	if !ebiten.IsFocused() {
		return nil
	}

	if len(p.animation.Parts) == 0 {
		return nil
	}

	// Play/Stop by space key.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.playing = !p.playing
	}

	// Arrow keys
	if !p.playing {
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			p.currentTick--
			if p.currentTick < 0 {
				p.currentTick = 0
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			p.currentTick++
			if p.currentTick >= p.maxTick {
				p.currentTick = p.maxTick - 1
			}
		}
	}

	// Cursor
	cursorX, cursorY := ebiten.CursorPosition()
	isCursorOnAnimation := cursorX >= p.offsetX && cursorX <= p.offsetX+p.width && cursorY >= p.offsetY && cursorY <= p.offsetY+p.height
	if !isCursorOnAnimation {
		return nil
	}
	xOnAnimation := cursorX - p.offsetX
	yOnAnimation := cursorY - p.offsetY
	isCurorOnPlayer := xOnAnimation >= p.barX && xOnAnimation <= p.barX+p.barWidth && yOnAnimation >= p.barY && yOnAnimation <= p.barY+p.barHeight
	if !isCurorOnPlayer {
		return nil
	}
	p.cursolOnBar = true
	ebiten.SetCursorShape(ebiten.CursorShapePointer)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		p.playing = false
		p.currentTick = int(float64(xOnAnimation-p.barX) / float64(p.barWidth) * float64(p.maxTick))
		if p.currentTick >= p.maxTick {
			p.currentTick = p.maxTick - 1
		}
		for i, index := range p.indexes {
			if p.currentTick >= index {
				p.currentPart = i
			} else {
				break
			}
		}
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	// Fill background.
	bg := ebiten.NewImage(screen.Bounds().Dx()-constant.MenuWidth, p.height)
	bgOp := &ebiten.DrawImageOptions{}
	bgOp.GeoM.Translate(float64(constant.MenuWidth), float64(p.offsetY))

	// Play/Stop button.
	for _, b := range p.buttons {
		defer b.Draw(screen)
	}
	for _, cl := range p.customLinks {
		if cl.id == "# Properties" && len(p.animation.Parts) == 0 {
			break
		}
		defer cl.Draw(screen)
	}

	// Final Draw.
	defer screen.DrawImage(bg, bgOp)

	// Animation Frame.
	frameLine := ebiten.NewImage(p.animation.Width+2, p.animation.Height+2)
	frameLine.Fill(color.Black)
	frame := ebiten.NewImage(p.animation.Width, p.animation.Height)
	frame.Fill(color.White)
	{
		isGray := false
		boxSize := constant.DefaultAnimationSize / 5
		maxX := int(math.Ceil(float64(p.animation.Width) / float64(boxSize)))
		maxY := int(math.Ceil(float64(p.animation.Height) / float64(boxSize)))
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
	x := float64(bg.Bounds().Dx()/2) - float64(p.animation.Width)/2
	y := float64(bg.Bounds().Dy()/3) - float64(p.animation.Height)/2
	frameOp := &ebiten.DrawImageOptions{}
	frameLineOp := &ebiten.DrawImageOptions{}
	frameOp.GeoM.Translate(float64(x), float64(y))
	frameLineOp.GeoM.Translate(float64(x-1), float64(y-1))
	bg.DrawImage(frameLine, frameLineOp)

	// Draw part.
	if p.currentPart >= 0 {
		part := p.animation.Parts[p.currentPart]
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
			op.GeoM.Translate(float64(p.animation.Width/2), float64(p.animation.Height/2))
			frame.DrawImage(img, op)
		}
	}

	bg.DrawImage(frame, frameOp)

	// Player.
	player := ebiten.NewImage(p.barWidth, p.barHeight)
	player.Fill(color.Gray{Y: constant.PlayerBarGrayY})
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(10, float64(p.barY))
	prs := int(float64(p.currentTick) / float64(p.maxTick) * float64(player.Bounds().Dx()))
	progress := ebiten.NewImage(2, player.Bounds().Dy())
	progress.Fill(color.Black)
	progressOp := &ebiten.DrawImageOptions{}
	progressOp.GeoM.Translate(float64(prs), 0)
	// draw index positions
	indexPos := ebiten.NewImage(2, player.Bounds().Dy())
	indexPos.Fill(color.Gray{Y: constant.PlayerBarIndexGrayY})
	for i, index := range p.indexes {
		x := int(float64(index) / float64(p.maxTick) * float64(player.Bounds().Dx()))
		indexPosOp := &ebiten.DrawImageOptions{}
		indexPosOp.GeoM.Translate(float64(x), 0)
		player.DrawImage(indexPos, indexPosOp)
		if i == p.currentPart {
			// Fill current part.
			width := int(math.Ceil(float64(p.animation.Parts[i].Length) / float64(p.maxTick) * float64(player.Bounds().Dx())))
			currentPart := ebiten.NewImage(width, player.Bounds().Dy())
			currentPart.Fill(color.Gray{Y: constant.PlayerBarCurrentGrayY})
			currentPartOp := &ebiten.DrawImageOptions{}
			currentPartOp.GeoM.Translate(float64(x), 0)
			player.DrawImage(currentPart, currentPartOp)
		}
	}
	player.DrawImage(progress, progressOp)
	bg.DrawImage(player, playerOp)
	s := fmt.Sprintf("Part: %d / %d\nTick: %d / %d\nSec: %0.2f / %0.2f", p.currentPart+1, len(p.animation.Parts), p.currentTick, p.maxTick, float64(p.currentTick)/float64(ebiten.TPS()), float64(p.maxTick)/float64(ebiten.TPS()))
	bs := text.BoundString(p.font, s)
	text.Draw(bg, s, p.font, 10, p.barY-bs.Dy(), color.Black)
}

func (p *Player) Layout(outsideWidth, outsideHeight int) {
	// Resize
	p.height = outsideHeight - outsideHeight/3
	p.offsetY = outsideHeight / 3
	p.offsetX = constant.MenuWidth
	p.width = outsideWidth - p.offsetX
	p.barWidth = p.width - 20
	p.barHeight = 30
	p.barX = 10
	p.barY = p.height - p.barHeight - 10

	totalWidth := 0
	for _, b := range p.buttons {
		totalWidth += b.width + 5
	}
	startX := p.offsetX + (p.width-totalWidth)/2 + 5
	for _, b := range p.buttons {
		b.MoveTo(startX, p.offsetY+p.height-75)
		startX += b.width + 5
	}
	defaultBs := text.BoundString(p.font, "DEFAULT")
	for i, l := range p.customLinks {
		l.MoveTo(p.offsetX+5, p.offsetY+5+((defaultBs.Dy()+5)*i))
	}
}

func (p *Player) CanExport() bool {
	return len(p.animation.Parts) > 0
}

func (p *Player) RawAnimation() *animation.Animation {
	return p.animation
}

func (p *Player) Stop() {
	p.playing = false
}

func (p *Player) Import(a *animation.Animation) {
	p.animation = a
	p.resetIndexes()
}

package ui

import (
	"image/color"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type Button struct {
	x, y, width, height int
	label               string
	cursorOn            bool
	font                font.Face
	disabled            bool
	onClick             func()
}

func NewButton(x, y, width, height int, label string, onClick func()) *Button {
	tt, _ := opentype.Parse(goregular.TTF)
	font, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})
	return &Button{
		x:        x,
		y:        y,
		width:    width,
		height:   height,
		label:    label,
		font:     font,
		disabled: false,
		onClick:  onClick,
	}
}

func (b *Button) Update() error {
	if !ebiten.IsFocused() || b.disabled {
		b.cursorOn = false
		return nil
	}
	cursorX, cursorY := ebiten.CursorPosition()
	if cursorX >= b.x && cursorX <= b.x+b.width && cursorY >= b.y && cursorY <= b.y+b.height {
		b.cursorOn = true
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		if b.cursorOn {
			b.cursorOn = false
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
	}
	if b.cursorOn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.onClick()
	}

	return nil
}

func (b *Button) Draw(screen *ebiten.Image) {
	if b.cursorOn {
		img := ebiten.NewImage(b.width, b.height)
		img.Fill(color.Gray{Y: constant.ButtonGrayY})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(b.x), float64(b.y))
		screen.DrawImage(img, op)
	}
	defaultBs := text.BoundString(b.font, "DEFAULT")
	bs := text.BoundString(b.font, b.label)
	var clr color.Color
	if b.disabled {
		clr = color.Gray{Y: constant.DisabledGrayY}
	} else {
		clr = color.Black
	}
	text.Draw(screen, b.label, b.font, b.x+(b.width-bs.Dx())/2, b.y+b.height-defaultBs.Dy(), clr)
}

func (b *Button) Layout(outsideWidth, outsideHeight int) {
}

func (b *Button) MoveTo(x, y int) {
	b.x = x
	b.y = y
}

func (b *Button) Label() string {
	return b.label
}

func (b *Button) SetLabel(label string) {
	b.label = label
}

func (b *Button) SetDisabled(disabled bool) {
	b.disabled = disabled
}

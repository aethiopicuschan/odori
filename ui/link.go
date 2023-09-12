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

type Link struct {
	x, y, width, height int
	id                  string
	label               string
	cursorOn            bool
	font                font.Face
	disabled            bool
	onClick             func()
}

func NewLink(x, y int, label string, onClick func()) *Link {
	tt, _ := opentype.Parse(goregular.TTF)
	font, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})

	return &Link{
		x:        x,
		y:        y,
		width:    0,
		height:   0,
		id:       label,
		label:    label,
		cursorOn: false,
		font:     font,
		disabled: false,
		onClick:  onClick,
	}
}

func (l *Link) Update() error {
	bs := text.BoundString(l.font, "DEFAULT")
	l.height = bs.Dy()
	bs = text.BoundString(l.font, l.label)
	l.width = bs.Dx()

	if !ebiten.IsFocused() || l.onClick == nil || l.disabled {
		l.cursorOn = false
		return nil
	}
	cursorX, cursorY := ebiten.CursorPosition()
	if cursorX >= l.x && cursorX <= l.x+l.width && cursorY >= l.y && cursorY <= l.y+l.height {
		l.cursorOn = true
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		l.cursorOn = false
	}
	if l.cursorOn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		l.onClick()
	}

	return nil
}

func (l *Link) Draw(screen *ebiten.Image) {
	bs := text.BoundString(l.font, l.label)
	var clr color.Color
	if l.onClick == nil {
		clr = color.Black
	} else if l.disabled {
		clr = color.Gray{Y: constant.DisabledGrayY}
	} else {
		clr = color.RGBA{R: 26, G: 13, B: 171, A: 255}
	}
	text.Draw(screen, l.label, l.font, l.x+(l.width-bs.Dx())/2, l.y+l.height, clr)
	// draw line if cursor on
	if l.cursorOn {
		img := ebiten.NewImage(l.width, 1)
		img.Fill(color.RGBA{R: 26, G: 13, B: 171, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(l.x), float64(l.y+l.height+1))
		screen.DrawImage(img, op)
	}
}

func (l *Link) Layout(outsideWidth, outsideHeight int) {
}

func (l *Link) MoveTo(x, y int) {
	l.x = x
	l.y = y
}

func (l *Link) SetLabel(label string) {
	l.label = label
}

func (l *Link) SetDisabled(disabled bool) {
	l.disabled = disabled
}

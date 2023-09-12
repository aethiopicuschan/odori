package ui

import (
	"image/color"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type level int

const (
	INFO level = iota
	WARN
	ERROR
)

type notice struct {
	level   level
	message string
	count   int
}

type Noticer struct {
	noticeHeight int
	notices      []*notice
	font         font.Face
}

func NewNoticer(noticeHeight int) *Noticer {
	tt, _ := opentype.Parse(goregular.TTF)
	font, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 12,
		DPI:  72,
	})

	return &Noticer{
		noticeHeight: noticeHeight,
		notices:      []*notice{},
		font:         font,
	}
}

func (n *Noticer) AddNotice(level level, message string) {
	n.notices = append(n.notices, &notice{
		level:   level,
		message: message,
		count:   constant.NoticeTime * ebiten.TPS(),
	})
}

func (n *Noticer) Update() error {
	for _, notice := range n.notices {
		notice.count--
		if notice.count <= 0 {
			n.notices = n.notices[1:]
		}
	}
	return nil
}

func (n *Noticer) Draw(screen *ebiten.Image) {
	for i, notice := range n.notices {
		img := ebiten.NewImage(screen.Bounds().Dx()-20, n.noticeHeight)
		switch notice.level {
		case INFO:
			img.Fill(color.RGBA{R: 58, G: 110, B: 165, A: 255})
		case WARN:
			img.Fill(color.RGBA{R: 255, G: 128, B: 0, A: 255})
		case ERROR:
			img.Fill(color.RGBA{R: 204, G: 51, B: 0, A: 255})
		}
		op := &ebiten.DrawImageOptions{}
		x := 10
		y := screen.Bounds().Dy() - 10 - ((n.noticeHeight + 10) * (i + 1))
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)
		bs := text.BoundString(n.font, notice.message)
		text.Draw(screen, notice.message, n.font, x+(img.Bounds().Dx()-bs.Dx())/2, y+img.Bounds().Dy()-bs.Dy(), color.White)
	}
}

func (n *Noticer) Layout(outsideWidth, outsideHeight int) {
}

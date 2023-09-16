package animation

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"

	"github.com/aethiopicuschan/odori/constant"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soniakeys/quant/median"
)

type Animation struct {
	Parts  []*Part `json:"parts"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
}

func NewAnimation() *Animation {
	return &Animation{
		Parts:  []*Part{},
		Width:  constant.DefaultAnimationSize,
		Height: constant.DefaultAnimationSize,
	}
}

func (a *Animation) CanExport() bool {
	return len(a.Parts) > 0
}

func (a *Animation) ExportAsGif(path string) (err error) {
	if !a.CanExport() {
		return errors.New("can not export")
	}
	outGif := &gif.GIF{}
	for _, part := range a.Parts {
		frame := ebiten.NewImage(a.Width, a.Height)
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
			op.GeoM.Translate(float64(a.Width/2), float64(a.Height/2))
			frame.DrawImage(img, op)
		}
		// パレットを作成
		q := median.Quantizer(255)
		p := color.Palette{image.Transparent}
		p = append(p, q.Quantize(make(color.Palette, 0, 255), frame)...)
		// frameをPalettedに変換してoutGifに追加
		paletted := image.NewPaletted(image.Rect(0, 0, a.Width, a.Height), p)
		draw.Draw(paletted, paletted.Rect, frame, image.Point{0, 0}, draw.Src)
		outGif.Image = append(outGif.Image, paletted)
		// GifのDelayは1/100なので、変換する
		delay := float64(part.Length) / float64(ebiten.TPS()) * 100
		outGif.Delay = append(outGif.Delay, int(delay))
		outGif.Disposal = append(outGif.Disposal, gif.DisposalBackground)
	}
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	err = gif.EncodeAll(file, outGif)
	return
}

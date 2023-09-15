package sprite

import (
	"image"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image *ebiten.Image
	id    string
}

func NewSprite(img image.Image) (sprite Sprite) {
	return Sprite{
		Image: ebiten.NewImageFromImage(img),
		id:    uuid.NewString(),
	}
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func NewSpriteFromRects(img image.Image, rects []image.Rectangle) (sprites []Sprite) {
	for _, rect := range rects {
		sprites = append(sprites, NewSprite(img.(SubImager).SubImage(rect)))
	}
	return
}

func NewEmptySprite() (sprite Sprite) {
	return Sprite{
		Image: nil,
		id:    uuid.NewString(),
	}
}

func (s *Sprite) IsEmpty() bool {
	return s.Image == nil
}

func (s *Sprite) Id() string {
	return s.id
}

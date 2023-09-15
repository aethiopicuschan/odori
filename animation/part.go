package animation

import "github.com/aethiopicuschan/odori/sprite"

type Part struct {
	Sprite  sprite.Sprite
	Scale   float64
	DiffX   int
	DiffY   int
	Reverse bool
	Length  int
}

func NewPart(sprite sprite.Sprite, len int) *Part {
	return &Part{
		Sprite:  sprite,
		Scale:   1.0,
		DiffX:   0,
		DiffY:   0,
		Reverse: false,
		Length:  len,
	}
}

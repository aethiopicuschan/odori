package animation

import (
	"github.com/aethiopicuschan/odori/sprite"
)

type Part struct {
	Sprite  sprite.Sprite `json:"sprite"`
	Scale   float64       `json:"scale"`
	DiffX   int           `json:"diffX"`
	DiffY   int           `json:"diffY"`
	Reverse bool          `json:"reverse"`
	Length  int           `json:"length"`
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

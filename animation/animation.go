package animation

import (
	"image"

	"github.com/aethiopicuschan/odori/constant"
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

// 永続化用モデル
type AnimationP struct {
	Name        string                 `json:"name"`
	Animation   *Animation             `json:"animation"`
	SpriteSheet map[string]image.Point `json:"spriteSheet"`
}

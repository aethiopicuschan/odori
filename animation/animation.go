package animation

import "github.com/aethiopicuschan/odori/constant"

type Animation struct {
	Parts  []*Part
	Width  int
	Height int
}

func NewAnimation() *Animation {
	return &Animation{
		Parts:  []*Part{},
		Width:  constant.DefaultAnimationSize,
		Height: constant.DefaultAnimationSize,
	}
}

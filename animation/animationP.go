package animation

import "image"

// 永続化用モデル
type AnimationP struct {
	Name        string                     `json:"name"`
	Animation   *Animation                 `json:"animation"`
	SpriteSheet map[string]image.Rectangle `json:"spriteSheet"`
}

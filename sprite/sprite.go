package sprite

import (
	"encoding/json"
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

func NewSpriteWithId(img image.Image, id string) (sprite Sprite) {
	return Sprite{
		Image: ebiten.NewImageFromImage(img),
		id:    id,
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

func NewSpritesFromRectMap(img image.Image, rectMap map[string]image.Rectangle) (sprites []Sprite) {
	for id, rect := range rectMap {
		sprites = append(sprites, NewSpriteWithId(img.(SubImager).SubImage(rect), id))
	}
	return
}

func NewEmptySprite() (sprite Sprite) {
	return Sprite{
		Image: nil,
		id:    "",
	}
}

func (s *Sprite) IsEmpty() bool {
	return s.Image == nil
}

func (s *Sprite) Id() string {
	return s.id
}

// スプライトの永続化モデル
type SpriteP struct {
	Id      string `json:"id"`
	IsEmpty bool   `json:"isEmpty"`
}

func (s *Sprite) MarshalJSON() ([]byte, error) {
	id := s.id
	if s.IsEmpty() {
		id = ""
	}
	spriteP := SpriteP{
		Id:      id,
		IsEmpty: s.IsEmpty(),
	}
	return json.Marshal(spriteP)
}

func (s *Sprite) UnmarshalJSON(bytes []byte) (err error) {
	spriteP := SpriteP{}
	err = json.Unmarshal(bytes, &spriteP)
	if err != nil {
		return
	}
	if spriteP.IsEmpty {
		s.Image = nil
	} else {
		// 仮の画像を入れておく
		s.Image = ebiten.NewImage(1, 1)
		s.id = spriteP.Id
	}
	return
}

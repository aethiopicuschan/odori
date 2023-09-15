package io

import (
	"image"

	"github.com/aethiopicuschan/kaban/merge"
	"github.com/aethiopicuschan/odori/sprite"
)

type WriteSpriteSheetResult struct {
	PointsMap map[string]image.Point
	Err       error
}

func WriteSpriteSheet(ch chan WriteSpriteSheetResult, sprites []sprite.Sprite, path string) {
	result := WriteSpriteSheetResult{}
	defer func() {
		ch <- result
	}()
	imgs := make([]image.Image, len(sprites))
	for i, s := range sprites {
		if !s.IsEmpty() {
			imgs[i] = s.Image
		}
	}
	img, points, err := merge.Merge(imgs)
	if err != nil {
		result.Err = err
		return
	}
	err = writePng(img, path)
	if err != nil {
		result.Err = err
		return
	}
	result.PointsMap = make(map[string]image.Point)
	for i, point := range points {
		result.PointsMap[sprites[i].Id()] = point
	}
}
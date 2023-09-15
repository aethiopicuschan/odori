package io

import (
	"image"
	"os"

	"github.com/aethiopicuschan/kaban/merge"
	"github.com/aethiopicuschan/odori/sprite"
)

func Write(ch chan error, bytes []byte, path string) {
	file, err := os.Create(path)
	if err != nil {
		ch <- err
		close(ch)
		return
	}
	defer file.Close()
	_, err = file.Write(bytes)
	ch <- err
	close(ch)
}

type WriteSpriteSheetResult struct {
	PointsMap map[string]image.Point
	Err       error
}

func WriteSpriteSheet(ch chan WriteSpriteSheetResult, sprites []sprite.Sprite, path string) {
	result := WriteSpriteSheetResult{}
	defer func() {
		ch <- result
		close(ch)
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

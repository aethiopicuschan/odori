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
		return
	}
	defer file.Close()
	_, err = file.Write(bytes)
	ch <- err
}

type WriteSpriteSheetResult struct {
	RectsMap map[string]image.Rectangle
	Err      error
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
	img, rects, err := merge.Merge(imgs)
	if err != nil {
		result.Err = err
		return
	}
	err = writePng(img, path)
	if err != nil {
		result.Err = err
		return
	}
	result.RectsMap = make(map[string]image.Rectangle)
	for i, rect := range rects {
		result.RectsMap[sprites[i].Id()] = rect
	}
}

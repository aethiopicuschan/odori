package io

import (
	"github.com/aethiopicuschan/kaban/detection"
	"github.com/aethiopicuschan/odori/sprite"
)

type ReadResult struct {
	Sprite sprite.Sprite
	Err    error
}

func Read(ch chan ReadResult, path string) {
	result := ReadResult{}
	defer func() {
		ch <- result
	}()
	img, err := readPng(path)
	if err != nil {
		result.Err = err
	} else {
		result.Sprite = sprite.NewSprite(img)
	}
}

type ReadSpriteSheetResult struct {
	Sprites []sprite.Sprite
	Err     error
}

func ReadSpriteSheet(ch chan ReadSpriteSheetResult, path string) {
	result := ReadSpriteSheetResult{}
	defer func() {
		ch <- result
	}()
	img, err := readPng(path)
	if err != nil {
		result.Err = err
		return
	}
	rects, err := detection.Detect(img)
	if err != nil {
		result.Err = err
	} else {
		result.Sprites = sprite.NewSpriteFromRects(img, rects)
	}
}

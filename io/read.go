package io

import (
	"os"

	"github.com/aethiopicuschan/kaban/detection"
	"github.com/aethiopicuschan/odori/sprite"
)

type ReadResult struct {
	Bytes []byte
	Path  string
	Err   error
}

func Read(ch chan ReadResult, path string) {
	result := ReadResult{
		Path: path,
	}
	defer func() {
		ch <- result
	}()
	result.Bytes, result.Err = os.ReadFile(path)
}

type ReadSpriteResult struct {
	Sprite sprite.Sprite
	Path   string
	Err    error
}

func ReadSprite(ch chan ReadSpriteResult, path string) {
	result := ReadSpriteResult{
		Path: path,
	}
	defer func() {
		ch <- result
	}()
	img, err := ReadPng(path)
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
	img, err := ReadPng(path)
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

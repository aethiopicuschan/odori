package io

import (
	"image"
	_ "image/png"
	"os"
)

func readPng(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

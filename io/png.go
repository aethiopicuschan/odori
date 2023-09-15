package io

import (
	"image"
	"image/png"
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

func writePng(img image.Image, path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, img)
	return
}

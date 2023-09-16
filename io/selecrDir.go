package io

import (
	"github.com/ncruces/zenity"
)

type SelectDirResult struct {
	Path string
	Err  error
}

func SelectDir(ch chan SelectDirResult) {
	result := SelectDirResult{}
	defer func() {
		ch <- result
	}()

	result.Path, result.Err = zenity.SelectFile(zenity.Directory())
}

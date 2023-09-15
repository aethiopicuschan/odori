package io

import "github.com/ncruces/zenity"

type PickResult struct {
	Paths []string
	Err   error
}

func Pick(ch chan PickResult) {
	pickResult := PickResult{}
	defer func() {
		ch <- pickResult
		close(ch)
	}()
	path, err := zenity.SelectFile(
		zenity.FileFilters{
			{Name: "Image files", Patterns: []string{"*.png"}, CaseFold: true},
		})
	if err != nil {
		pickResult.Err = err
	} else {
		pickResult.Paths = []string{path}
	}
}

func PickMultiple(ch chan PickResult) {
	pickResult := PickResult{}
	defer func() {
		ch <- pickResult
		close(ch)
	}()
	paths, err := zenity.SelectFileMultiple(
		zenity.FileFilters{
			{Name: "Image files", Patterns: []string{"*.png"}, CaseFold: true},
		})
	if err != nil {
		pickResult.Err = err
	} else {
		pickResult.Paths = paths
	}
}

package io

import "github.com/ncruces/zenity"

type EntryResult struct {
	Input string
	Err   error
}

func Entry(ch chan EntryResult, title, text, def string) {
	result := EntryResult{}
	defer func() {
		ch <- result
		close(ch)
	}()
	result.Input, result.Err = zenity.Entry(text, zenity.EntryText(def), zenity.Title(title))
}

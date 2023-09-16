package io

import "github.com/ncruces/zenity"

type PickResult struct {
	Path string
	Err  error
}

type PickOption struct {
	Name     string
	Patterns []string
}

func WithName(name string) func(*PickOption) {
	return func(o *PickOption) {
		o.Name = name
	}
}

func WithPatterns(patterns []string) func(*PickOption) {
	return func(o *PickOption) {
		o.Patterns = patterns
	}
}

func Pick(ch chan PickResult, options ...func(*PickOption)) {
	pickResult := PickResult{}
	defer func() {
		ch <- pickResult
	}()
	opt := &PickOption{
		Name: "Select file",
	}
	for _, o := range options {
		o(opt)
	}

	path, err := zenity.SelectFile(
		zenity.FileFilters{
			{Name: opt.Name, Patterns: opt.Patterns, CaseFold: true},
		})
	if err != nil {
		pickResult.Err = err
	} else {
		pickResult.Path = path
	}
}

type PickMultipleResult struct {
	Paths []string
	Err   error
}

func PickMultiple(ch chan PickMultipleResult, options ...func(*PickOption)) {
	pickResult := PickMultipleResult{}
	defer func() {
		ch <- pickResult
	}()
	opt := &PickOption{
		Name: "Select files",
	}
	for _, o := range options {
		o(opt)
	}

	paths, err := zenity.SelectFileMultiple(
		zenity.FileFilters{
			{Name: opt.Name, Patterns: opt.Patterns, CaseFold: true},
		})
	if err != nil {
		pickResult.Err = err
	} else {
		pickResult.Paths = paths
	}
}

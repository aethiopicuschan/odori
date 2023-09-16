package io

import "github.com/ncruces/zenity"

type PickResult struct {
	Path string
	Err  error
}

type PickOption struct {
	toSave      bool
	defaultName string
	name        string
	patterns    []string
}

func WithName(name string) func(*PickOption) {
	return func(o *PickOption) {
		o.name = name
	}
}

func WithPatterns(patterns []string) func(*PickOption) {
	return func(o *PickOption) {
		o.patterns = patterns
	}
}

func WithToSave(defaultName string) func(*PickOption) {
	return func(o *PickOption) {
		o.toSave = true
		o.defaultName = defaultName
	}
}

func Pick(ch chan PickResult, options ...func(*PickOption)) {
	pickResult := PickResult{}
	defer func() {
		ch <- pickResult
	}()
	opt := &PickOption{
		name: "Select file",
	}
	for _, o := range options {
		o(opt)
	}

	if opt.toSave {
		pickResult.Path, pickResult.Err = zenity.SelectFileSave(
			zenity.ConfirmOverwrite(),
			zenity.Filename(opt.defaultName),
			zenity.FileFilters{
				{Name: opt.name, Patterns: opt.patterns, CaseFold: true},
			})
	} else {
		pickResult.Path, pickResult.Err = zenity.SelectFile(
			zenity.FileFilters{
				{Name: opt.name, Patterns: opt.patterns, CaseFold: true},
			})
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
		name: "Select files",
	}
	for _, o := range options {
		o(opt)
	}

	paths, err := zenity.SelectFileMultiple(
		zenity.FileFilters{
			{Name: opt.name, Patterns: opt.patterns, CaseFold: true},
		})
	if err != nil {
		pickResult.Err = err
	} else {
		pickResult.Paths = paths
	}
}

package io

import "github.com/ncruces/zenity"

type QuestionResult struct {
	Answer bool
}

func Question(ch chan QuestionResult, title, text string) {
	result := QuestionResult{}
	defer func() {
		ch <- result
	}()
	err := zenity.Question(text, zenity.Title(title))
	result.Answer = err == nil
}

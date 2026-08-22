package main

import (
	"regexp"
	"strings"
)

type Tokenizer struct {
	stopWords      map[string]bool
	minTokenLength int
}

func NewTokenizer(stopWords []string, minTokenLength int) *Tokenizer {
	if minTokenLength == 0 {
		minTokenLength = 2
	}
	sw := map[string]bool{
		"the": true, "is": true, "and": true, "or": true, "to": true,
		"in": true, "on": true, "for": true, "of": true, "a": true, "an": true,
	}
	for _, w := range stopWords {
		sw[w] = true
	}
	return &Tokenizer{stopWords: sw, minTokenLength: minTokenLength}
}

var punctRegex = regexp.MustCompile(`[^\p{L}\p{N}\s]+`)

func (t *Tokenizer) Tokenize(text string) []string {
	if text == "" {
		return nil
	}
	cleaned := punctRegex.ReplaceAllString(strings.ToLower(text), "")
	parts := strings.Fields(cleaned)
	var tokens []string
	for _, part := range parts {
		if len(part) >= t.minTokenLength && !t.stopWords[part] {
			tokens = append(tokens, part)
		}
	}
	return tokens
}

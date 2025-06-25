package main

import (
	"sort"
	"strings"
	"unicode"
)

func Top10(text string) []string {
	if text == "" {
		return []string{}
	}

	words := make(map[string]int)

	for _, val := range strings.Fields(text) {
		if bufStr := getWord(val); bufStr != "" {
			words[bufStr]++
		}
	}
	keys := make([]string, 0, len(words))
	for k := range words {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		valI, valJ := keys[i], keys[j]
		if words[valI] != words[valJ] {
			return words[valI] > words[valJ]
		}
		return valI < valJ
	})
	if len(keys) > 10 {
		keys = keys[:10]
	}

	return keys
}

func getWord(text string) string {
	if isPunctString(text) {
		return text
	}
	return strings.TrimFunc(strings.ToLower(text), func(r rune) bool {
		return unicode.IsDigit(r) || unicode.IsPunct(r)
	})
}

func isPunctString(text string) bool {
	runes := []rune(text)
	for _, v := range runes {
		if !unicode.IsPunct(v) {
			return false
		}
	}
	return len(runes) > 1
}

func main() {}

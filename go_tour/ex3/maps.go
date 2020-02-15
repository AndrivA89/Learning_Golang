package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

var match = false

func WordCount(s string) map[string]int {
	results := make(map[string]int)
	words := strings.Fields(s)

	for _, word := range words {
		for result, _ := range results {
			if result == word {
				match = true
			}
		}
		if match {
			results[word] += 1
		} else {
			results[word] = 1
		}
	}
	return results
}

func main() {
	wc.Test(WordCount)
}

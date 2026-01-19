package gen

import (
	"math/rand"
	"strconv"
)

func WithNumbers(words []string) []string {
	result := make([]string, len(words))
	copy(result, words)

	for i, b := range Weighted(len(words), map[bool]int{
		false: 93,
		true:  7,
	}) {
		if b {
			result[i] = strconv.Itoa(rand.Intn(9999) + 1)
		}
	}

	return result
}

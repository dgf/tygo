package gen

import (
	"math/rand"
	"slices"
	"strconv"
)

func WithNumbers(weight int, words []string) []string {
	result := slices.Clone(words)

	for i, b := range SampleWeightedDist(len(words), map[bool]int{
		false: 100 - weight,
		true:  weight,
	}) {
		if b {
			result[i] = strconv.Itoa(rand.Intn(9999) + 1)
		}
	}

	return result
}

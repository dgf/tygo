package gen

import (
	"slices"
	"strconv"
)

const MaxRandomNumber = 9999

func WithNumbers(weight int, words []string) []string {
	result := slices.Clone(words)
	count := len(words)
	dist := map[bool]int{false: 100 - weight, true: weight}

	for i, b := range SampleWeightedDist(count, dist) {
		if b {
			result[i] = strconv.Itoa(randGen.Intn(MaxRandomNumber) + 1)
		}
	}

	return result
}

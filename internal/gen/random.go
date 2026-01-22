package gen

import (
	"math/rand"
	"sort"
	"time"
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func Weighted[E comparable](amount int, distributions map[E]int) []E {
	type distSum struct {
		dist E
		sum  int64
	}

	count := len(distributions)
	result := make([]E, amount)
	distSums := make([]distSum, count)

	c := 0
	sum := int64(0)

	for k, v := range distributions {
		sum += int64(v)
		distSums[c] = distSum{k, sum}
		c++
	}

	for a := range amount {
		n := randGen.Int63n(sum) + 1

		s := sort.Search(count, func(c int) bool {
			return distSums[c].sum >= n
		})

		result[a] = distSums[s].dist
	}

	return result
}

func WeightedRandomList(amount int, words []string) []string {
	count := len(words)
	dists := make(map[string]int, count)

	for a, w := range words {
		dists[w] = count - a
	}

	return Weighted(amount, dists)
}

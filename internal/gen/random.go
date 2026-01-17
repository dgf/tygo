package gen

import (
	"math/rand"
	"sort"
	"time"
)

var randomgen = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(n int) int {
	return randomgen.Intn(n)
}

func WeightedRandomList(amount int, words []string) []string {
	result := make([]string, amount)
	count := len(words)
	sum := (count * (count + 1)) / 2

	for a := range amount {
		n := RandomInt(sum)

		calls := 0
		s := sort.Search(count, func(c int) bool {
			calls++

			return ((c+1)*(c+2))/2 >= n
		})

		result[a] = words[count-s-1]
	}

	return result
}

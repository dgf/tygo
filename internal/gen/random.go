// Package gen provides word list generation with weighted sampling and punctuation.
package gen

import (
	"math/rand"
	"slices"
	"sort"
	"time"
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

type distRange[E comparable] struct {
	value   E
	prevCum int64
	weight  int
	cumSum  int64
}

func mapRanges[E comparable](dists map[E]int) (int64, []distRange[E]) {
	count := len(dists)
	ranges := make([]distRange[E], count)

	idx := 0
	cumSum := int64(0)

	for value, weight := range dists {
		prevCum := cumSum
		cumSum += int64(weight)
		ranges[idx] = distRange[E]{value, prevCum, weight, cumSum}
		idx++
	}

	return cumSum, ranges
}

func SampleWeighted[E comparable](count, windowSize int, dists map[E]int) []E {
	result := make([]E, count)
	sum, ranges := mapRanges(dists)

	if windowSize >= len(ranges) {
		windowSize = min(5, len(ranges)/2)
	}

	ridx := 0
	recent := make([]distRange[E], windowSize)

	for i := range count {
		var (
			recentWeight int64
			recentSorted []distRange[E]
		)

		if windowSize > 0 {
			for _, r := range recent {
				recentWeight += int64(r.weight)
			}

			recentSorted = slices.Clone(recent)
			sort.Slice(recentSorted, func(i int, j int) bool {
				return recentSorted[i].prevCum < recentSorted[j].prevCum
			})
		}

		n := randGen.Int63n(sum-recentWeight) + 1

		if windowSize > 0 {
			for _, r := range recentSorted {
				if n > r.prevCum {
					n += int64(r.weight)
				}
			}
		}

		s := sort.Search(len(ranges), func(idx int) bool {
			return ranges[idx].cumSum >= n
		})

		selected := ranges[s]
		result[i] = selected.value

		if windowSize > 0 {
			recent[ridx] = selected
			ridx = (ridx + 1) % windowSize
		}
	}

	return result
}

func SampleWeightedList(count, noRepeatWindow int, words []string) []string {
	weight := len(words)
	dists := make(map[string]int, weight)

	for i, w := range words {
		dists[w] = weight - i
	}

	return SampleWeighted(count, noRepeatWindow, dists)
}

func SampleWeightedDist[E comparable](count int, dist map[E]int) []E {
	return SampleWeighted(count, 0, dist)
}

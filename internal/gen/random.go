// Package gen provides word list generation with weighted sampling and punctuation.
package gen

import (
	"math/rand"
	"slices"
	"sort"
	"time"
)

const DefaultMinNoRepeatWindowSize = 5

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

type distRange[E comparable] struct {
	value   E
	prevCum int64
	weight  int
	cumSum  int64
}

func MapDistRanges[E comparable](dists map[E]int) (int64, []distRange[E]) {
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

func WeightRecent[E comparable](recent []distRange[E]) int64 {
	var weight int64
	for _, r := range recent {
		weight += int64(r.weight)
	}

	return weight
}

func AddRecentWeightOffset[E comparable](number int64, recent []distRange[E]) int64 {
	if len(recent) == 0 {
		return number
	}

	sorted := slices.Clone(recent)
	sort.Slice(sorted, func(i int, j int) bool {
		return sorted[i].prevCum < sorted[j].prevCum
	})

	for _, r := range sorted {
		if number > r.prevCum {
			number += int64(r.weight)
		}
	}

	return number
}

func SampleWeighted[E comparable](count, noRepeatWindow int, dists map[E]int) []E {
	result := make([]E, count)
	sum, ranges := MapDistRanges(dists)

	if noRepeatWindow >= len(ranges) {
		noRepeatWindow = min(DefaultMinNoRepeatWindowSize, len(ranges)/2)
	}

	ridx := 0
	recent := make([]distRange[E], noRepeatWindow)

	for i := range count {
		n := randGen.Int63n(sum-WeightRecent(recent)) + 1
		n = AddRecentWeightOffset(n, recent)

		s := sort.Search(len(ranges), func(idx int) bool {
			return ranges[idx].cumSum >= n
		})

		selected := ranges[s]
		result[i] = selected.value

		if noRepeatWindow > 0 {
			recent[ridx] = selected
			ridx = (ridx + 1) % noRepeatWindow
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

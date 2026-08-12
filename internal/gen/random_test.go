package gen_test

import (
	"testing"

	"github.com/dgf/tygo/internal/gen"
)

func TestSampleWeightedDist(t *testing.T) {
	t.Parallel()

	count := 100
	dist := map[string]int{"foo": 7, "bar": 3}
	list := gen.SampleWeightedDist(count, dist)

	if count != len(list) {
		t.Fatalf("expected %d results, got: %d", count, len(list))
	}

	counts := make(map[string]int, len(dist))
	for _, a := range list {
		counts[a]++
	}

	if counts["foo"] < count/2 {
		t.Errorf("expected more than %d of foo, got %d", count/2, counts["foo"])
	}

	if counts["bar"] > count/2 {
		t.Errorf("expected less than %d of bar, got %d", count/2, counts["bar"])
	}
}

func TestSampleWeightedList(t *testing.T) {
	t.Parallel()

	count := 1000
	words := []string{"foo", "bar"}
	list := gen.SampleWeightedList(count, 0, words)

	if count != len(list) {
		t.Fatalf("expected %d results, got: %d", count, len(list))
	}

	counts := make(map[string]int, len(words))
	for _, a := range list {
		counts[a]++
	}

	if counts["foo"] < count/2 {
		t.Errorf("expected more than %d of foo, got %d", count/2, counts["foo"])
	}

	if counts["bar"] > count/2 {
		t.Errorf("expected less than %d of bar, got %d", count/2, counts["bar"])
	}
}

func TestSampleWeightedList_UniqLenMinusOne(t *testing.T) {
	t.Parallel()

	count := 100
	words := []string{"one", "two", "foo", "bar", "baz"}
	list := gen.SampleWeightedList(count, len(words)-1, words)

	if count != len(list) {
		t.Fatalf("expected %d results, got: %d", count, len(list))
	}

	wordCounts := make(map[string]int, len(words))
	for _, w := range list {
		wordCounts[w]++
	}

	s := count / len(words)
	for _, w := range words {
		if wordCounts[w] != s {
			t.Errorf("expected %d results for %q, got: %d", s, w, wordCounts[w])
		}
	}
}

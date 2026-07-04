package gen_test

import (
	"testing"

	"github.com/dgf/tygo/internal/gen"
)

func TestSampleWeightedDist(t *testing.T) {
	t.Parallel()

	c := 100
	d := map[string]int{"foo": 7, "bar": 3}

	r := gen.SampleWeightedDist(c, d)

	if c != len(r) {
		t.Fatalf("expected %d results, got: %d", c, len(r))
	}

	counts := make(map[string]int, len(d))
	for _, a := range r {
		counts[a]++
	}

	if counts["foo"] < c/2 {
		t.Errorf("expected more than %d of foo, got %d", c/2, counts["foo"])
	}

	if counts["bar"] > c/2 {
		t.Errorf("expected less than %d of bar, got %d", c/2, counts["bar"])
	}
}

func TestSampleWeightedList(t *testing.T) {
	t.Parallel()

	c := 1000
	d := []string{"foo", "bar"}

	r := gen.SampleWeightedList(c, 0, d)

	if c != len(r) {
		t.Fatalf("expected %d results, got: %d", c, len(r))
	}

	counts := make(map[string]int, len(d))
	for _, a := range r {
		counts[a]++
	}

	if counts["foo"] < c/2 {
		t.Errorf("expected more than %d of foo, got %d", c/2, counts["foo"])
	}

	if counts["bar"] > c/2 {
		t.Errorf("expected less than %d of bar, got %d", c/2, counts["bar"])
	}
}

func TestSampleWeightedList_UniqLenMinusOne(t *testing.T) {
	t.Parallel()

	c := 100
	d := []string{"one", "two", "foo", "bar", "baz"}

	r := gen.SampleWeightedList(c, len(d)-1, d)

	if c != len(r) {
		t.Fatalf("expected %d results, got: %d", c, len(r))
	}

	wordCounts := make(map[string]int, len(d))
	for _, w := range r {
		wordCounts[w]++
	}

	s := c / len(d)
	for _, w := range d {
		if wordCounts[w] != s {
			t.Errorf("expected %d results for %q, got: %d", s, w, wordCounts[w])
		}
	}
}

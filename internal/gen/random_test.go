package gen_test

import (
	"testing"

	"github.com/dgf/tygo/internal/gen"
)

func TestWeighted(t *testing.T) {
	t.Parallel()

	c := 100
	d := map[string]int{"foo": 7, "bar": 3}

	r := gen.Weighted(c, d)

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

func TestWeightedList(t *testing.T) {
	t.Parallel()

	c := 1000
	d := []string{"foo", "bar"}

	r := gen.WeightedRandomList(c, d)

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

package test_test

import (
	"reflect"
	"testing"

	"github.com/dgf/tygo/internal/test"
)

func TestToLines(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name  string
		runes int
		words []string
		lines []test.Line
	}

	for _, tc := range []testCase{
		{
			"empty", 1,
			[]string{},
			[]test.Line{},
		},
		{
			"one", 3,
			[]string{"one"},
			[]test.Line{{{'o', 'n', 'e'}}},
		},
		{
			"two", 3,
			[]string{"one", "two"},
			[]test.Line{{{'o', 'n', 'e'}}, {{'t', 'w', 'o'}}},
		},
		{
			"one two", 7,
			[]string{"one", "two"},
			[]test.Line{{{'o', 'n', 'e'}, {'t', 'w', 'o'}}},
		},
		{
			"one two three", 12,
			[]string{"one", "two", "three"},
			[]test.Line{{{'o', 'n', 'e'}, {'t', 'w', 'o'}}, {{'t', 'h', 'r', 'e', 'e'}}},
		},
		{
			"go", 8, // FIXME use rune length! should and could fit in 6 = 1(len) + 1(space) + 4(len)
			[]string{"", "äöüß", "☠"},
			[]test.Line{{{''}, {'ä', 'ö', 'ü', 'ß'}}, {{'☠'}}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := test.ToLines(tc.runes, tc.words)
			if !reflect.DeepEqual(tc.lines, actual) {
				t.Errorf("invalid line transform\nwant:\n%U\ngot:\n%U\n", tc.lines, actual)
			}
		})
	}
}

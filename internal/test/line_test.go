package test_test

import (
	"reflect"
	"testing"

	"github.com/dgf/tygo/internal/test"
)

func TestToLines(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name  string
		runes int
		words []string
		lines []test.Line
	}{
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
			"go", 6,
			[]string{"", "äöüß", "☠"},
			[]test.Line{{{''}, {'ä', 'ö', 'ü', 'ß'}}, {{'☠'}}},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := test.ToLines(testCase.runes, testCase.words)
			if !reflect.DeepEqual(testCase.lines, actual) {
				t.Errorf("invalid line transform\nwant:\n%U\ngot:\n%U\n", testCase.lines, actual)
			}
		})
	}
}

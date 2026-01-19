package gen

import (
	"slices"
	"unicode"
)

type Punctuation int

const (
	None Punctuation = iota
	Period
	Comma
	Quotation
	Question
	Exclamation
	Brackets
	Braces
	Parenthesis
	Colon
	Semicolon
)

// Apostrophe
// Dash
// Hyphen
// Ellipsis

func PunctuationMarks(words []string) []string {
	result := make([]string, len(words))
	copy(result, words)

	// upper first word
	first := []rune(result[0])
	first[0] = unicode.ToTitle(first[0])
	result[0] = string(first)

	// last closed
	lastPunct := Weighted(1, map[Punctuation]int{
		Period:      5,
		Question:    4,
		Exclamation: 3,
	})
	result[len(result)-1] = applyPunctuation(lastPunct[0], result[len(result)-1])

	// apply random to all between
	for p, punct := range Weighted(len(words)-2, map[Punctuation]int{
		None:        84,
		Period:      12,
		Comma:       8,
		Quotation:   3,
		Question:    4,
		Exclamation: 3,
		Brackets:    2,
		Braces:      2,
		Parenthesis: 3,
		Colon:       3,
		Semicolon:   2,
	}) {
		result[p+1] = applyPunctuation(punct, result[p+1])

		if slices.Contains([]Punctuation{Period, Question, Exclamation}, punct) {
			r := []rune(result[p+2])
			r[0] = unicode.ToTitle(r[0])
			result[p+2] = string(r)
		}
	}

	return result
}

func applyPunctuation(punct Punctuation, word string) string {
	switch punct {
	case None:
		return word
	case Period:
		return word + "."
	case Comma:
		return word + ","
	case Quotation:
		return "\"" + word + "\""
	case Question:
		return word + "?"
	case Exclamation:
		return word + "!"
	case Brackets:
		return "[" + word + "]"
	case Braces:
		return "{" + word + "}"
	case Parenthesis:
		return "(" + word + ")"
	case Colon:
		return word + ":"
	case Semicolon:
		return word + ";"
	default:
		return word
	}
}

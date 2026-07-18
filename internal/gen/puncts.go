package gen

import (
	"slices"
	"unicode"
)

type Punctuation int

const (
	Word Punctuation = iota
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

func PunctuationMarks(words []string, dist map[Punctuation]int) []string {
	result := slices.Clone(words)

	// Uppercase first word
	first := []rune(result[0])
	first[0] = unicode.ToTitle(first[0])
	result[0] = string(first)

	// last closed
	lastPunct := SampleWeightedDist(1, map[Punctuation]int{
		Period:      dist[Period],
		Question:    dist[Question],
		Exclamation: dist[Exclamation],
	})
	result[len(result)-1] = applyPunctuation(lastPunct[0], result[len(result)-1])

	// apply random to all between
	for p, punct := range SampleWeightedDist(len(words)-2, dist) {
		result[p+1] = applyPunctuation(punct, result[p+1])

		if slices.Contains([]Punctuation{Period, Question, Exclamation}, punct) {
			r := []rune(result[p+2])
			r[0] = unicode.ToTitle(r[0])
			result[p+2] = string(r)
		}
	}

	return result
}

type Applicator func(word string) string

func Append(suffix string) Applicator {
	return func(word string) string {
		return word + suffix
	}
}

func Echo(word string) string {
	return word
}

func Enframe(prefix, suffix string) Applicator {
	return func(word string) string {
		return prefix + word + suffix
	}
}

var PunctuationApplicators = map[Punctuation]Applicator{
	Word:        Echo,
	Period:      Append("."),
	Comma:       Append(","),
	Quotation:   Enframe("\"", "\""),
	Question:    Append("?"),
	Exclamation: Append("!"),
	Brackets:    Enframe("[", "]"),
	Braces:      Enframe("{", "}"),
	Parenthesis: Enframe("(", ")"),
	Colon:       Append(":"),
	Semicolon:   Append(";"),
}

func applyPunctuation(punct Punctuation, word string) string {
	apply, ok := PunctuationApplicators[punct]

	if !ok {
		return word
	}

	return apply(word)
}

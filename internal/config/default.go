package config

func Default() Config {
	return Config{
		Dictionary:  "english",
		TopWords:    100,
		WordCount:   20,
		Width:       50,
		Numbers:     false,
		Punctuation: true,
		Distribution: Distribution{
			Word:        85,
			Number:      7,
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
		},
	}
}

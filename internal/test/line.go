package test

type Line [][]rune

func ToLines(cols int, words []string) []Line {
	lines := []Line{}
	line := Line{}
	lc := 0

	for _, word := range words {
		runes := []rune(word)
		if cols < lc+len(runes) {
			lines = append(lines, line)
			line = Line{}
			lc = 0
		}

		line = append(line, runes)
		lc += len(runes) + 1
	}

	if lc > 0 {
		lines = append(lines, line)
	}

	return lines
}

package test

type Grid [][]*Cell

func ToGrid(cols int, words []string) Grid {
	lines := ToLines(cols, words)
	grid := make(Grid, len(lines))

	for l, line := range lines {
		lcs := []*Cell{}

		for w, word := range line {
			for _, r := range word {
				lcs = append(lcs, Enqueue(r))
			}

			if w < len(line)-1 || l < len(grid)-1 {
				lcs = append(lcs, Enqueue(' '))
			}
		}

		grid[l] = lcs
	}

	return grid
}

package display

import (
	"fmt"
	"io"

	"github.com/dgf/tygo/internal/test"
)

func PrintGrid(out io.Writer, grid test.Grid) {
	for _, row := range grid {
		for _, cell := range row {
			PrintCell(out, cell)
		}

		NewLine(out)
	}

	// reset cursor to start
	CursorUp(out, len(grid))
	_, _ = fmt.Fprint(out, "\r")
}

func ResetGrid(out io.Writer, grid test.Grid) {
	_, _ = fmt.Fprint(out, "\r")
	_, _ = fmt.Fprint(out, EraseRightBelow)

	PrintGrid(out, grid)
}

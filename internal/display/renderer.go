package display

import (
	"io"

	"github.com/dgf/tygo/internal/test"
)

type Renderer struct {
	out  io.Writer
	row  int
	rows int
}

func NewRenderer(out io.Writer) *Renderer {
	return &Renderer{out: out, row: 0, rows: 0}
}

func (r *Renderer) Advance(cell *test.Cell, lineBreak bool) {
	if cell != nil {
		PrintCell(r.out, cell)
	}

	if lineBreak {
		r.row++
		NewLine(r.out)
	}
}

func (r *Renderer) Exit() {
	UndoLine(r.out)
}

func (r *Renderer) Next(grid test.Grid) {
	UndoLine(r.out)
	PrintLine(r.out, "---")
	PrintGrid(r.out, grid)

	r.row = 0
	r.rows = len(grid)
}

func (r *Renderer) Print(result test.Result) {
	skip := r.rows - r.row
	if skip > 1 {
		CursorDown(r.out, skip-1)
	}

	PrintResult(r.out, result)
}

func (r *Renderer) Reset(grid test.Grid) {
	ResetGrid(r.out, r.row)
	PrintGrid(r.out, grid)

	r.row = 0
	r.rows = len(grid)
}

func (r *Renderer) Retract(cells test.Cells) {
	CursorBack(r.out, len(cells)-1)

	for _, c := range cells {
		PrintCell(r.out, c)
	}

	CursorBack(r.out, len(cells))
}

package game

import "github.com/dgf/tygo/internal/test"

type Renderer interface {
	Advance(cell *test.Cell, lineBreak bool)
	Exit()
	Next(grid test.Grid)
	Print(result test.Result)
	Reset(grid test.Grid)
	Retract(cells test.Cells)
}

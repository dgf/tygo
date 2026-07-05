package display

import (
	"fmt"
	"io"
)

func NewLine(out io.Writer) {
	_, _ = fmt.Fprint(out, "\r\n")
}

func PrintLine(out io.Writer, line string) {
	_, _ = fmt.Fprint(out, line)
	NewLine(out)
}

func UndoLine(out io.Writer) {
	CursorUp(out, 1)
	_, _ = fmt.Fprint(out, EraseLineToEnd)
}

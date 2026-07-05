package display

import (
	"fmt"
	"io"
	"strconv"
)

func CursorBack(out io.Writer, n int) {
	_, _ = fmt.Fprint(out, CSI+strconv.Itoa(n)+"D")
}

func CursorDown(out io.Writer, n int) {
	_, _ = fmt.Fprint(out, CSI+strconv.Itoa(n)+"B")
}

func CursorUp(out io.Writer, n int) {
	_, _ = fmt.Fprint(out, CSI+strconv.Itoa(n)+"A")
}

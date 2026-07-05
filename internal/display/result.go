package display

import (
	"fmt"
	"io"

	"github.com/dgf/tygo/internal/test"
)

func PrintResult(out io.Writer, result test.Result) {
	NewLine(out)
	NewLine(out)
	fmt.Fprintf(out, "Result: %s", result)
	NewLine(out)
	NewLine(out)
	fmt.Fprint(out, "[ENTER] next or [ESC] to quit")
	NewLine(out)
}

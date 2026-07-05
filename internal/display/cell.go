package display

import (
	"fmt"
	"io"

	"github.com/dgf/tygo/internal/test"
)

func ColorCSI(s test.Status) string {
	switch s {
	case test.Active:
		return StyleActive
	case test.Failed:
		return StyleFailed
	case test.Passed:
		return StylePassed
	case test.Queued:
		return ""
	default:
		return ""
	}
}

func PrintCell(out io.Writer, c *test.Cell) {
	r := c.Rune

	if r == ' ' && c.Status == test.Failed {
		r = '_'
	}

	_, _ = fmt.Fprint(out, ColorCSI(c.Status)+string(r)+Reset)
}

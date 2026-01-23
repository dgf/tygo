package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/dgf/tygo/internal/dict"
	"github.com/dgf/tygo/internal/gen"
	"github.com/dgf/tygo/internal/test"
	"golang.org/x/term"
)

const (
	CSI           = "\033["
	Reset         = CSI + "0m"
	RestoreCursor = CSI + "u"
	SaveCursor    = CSI + "s"
	StyleActive   = CSI + "7m"
	StyleFailed   = CSI + "38;5;197m"
	StylePassed   = CSI + "2m"
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

func PrintGrid(out io.Writer, grid test.Grid) {
	for _, row := range grid {
		for _, cell := range row {
			PrintCell(out, cell)
		}

		_, _ = fmt.Fprint(out, "\r\n")
	}

	// reset cursor to start
	_, _ = fmt.Fprint(out, CSI+strconv.Itoa(len(grid))+"A\r")
}

var (
	dictParam  string
	dictionary dict.Dictionary
	dictTop    int

	fileName string

	wordCount   int
	termCols    int
	numbers     bool
	punctuation bool
)

func NextGrid(words []string) test.Grid {
	list := gen.WeightedRandomList(wordCount, words)

	if numbers {
		list = gen.WithNumbers(list)
	}

	if punctuation {
		list = gen.PunctuationMarks(list)
	}

	return test.ToGrid(termCols-1, list)
}

func Dictionary(name string) dict.Dictionary {
	switch name {
	case "german":
		return dict.German10K
	case "english":
		return dict.English10K
	default:
		return dict.English10K
	}
}

func init() {
	flag.StringVar(&dictParam, "dict", "english", "dictionary to use, available: german, english")
	flag.IntVar(&dictTop, "top", 100, "top count of words to load from source (dict or file)")

	flag.StringVar(&fileName, "file", "", "vocabulary JSON file with 'words' list")

	flag.IntVar(&wordCount, "count", 20, "number of words to include in the typing test")
	flag.IntVar(&termCols, "width", 50, "display width for the typing text")
	flag.BoolVar(&numbers, "nums", false, "enable number mode")
	flag.BoolVar(&punctuation, "punct", false, "enable punctuation marks")
}

func main() {
	flag.Parse()

	dictionary = Dictionary(dictParam)

	in := os.Stdin
	inFd := int(in.Fd())
	out := os.Stdout

	if !term.IsTerminal(inFd) {
		fmt.Fprintln(os.Stderr, "Use a termnial (requires a TTY)")
		os.Exit(1)

		return
	}

	termState, err := term.MakeRaw(inFd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Raw mode activation failed: %v\n", err)
		os.Exit(2)

		return
	}

	defer func() {
		_ = term.Restore(inFd, termState)

		if r := recover(); r != nil {
			fmt.Fprintf(out, "%v - %s", r, debug.Stack())
		}

		os.Exit(0)
	}()

	var words []string
	if len(fileName) > 0 {
		words = dict.LoadFile(fileName)
	} else {
		words = dict.LoadDict(dictionary, dictTop)
	}

	row := 0
	col := 0
	done := false

	grid := NextGrid(words)
	PrintGrid(out, grid)

	var startTime time.Time

	for {
		buf := make([]byte, 4)

		n, err := in.Read(buf)
		if err != nil {
			return // stdin closed > time to leave
		}

		// Ctrl+C (3) or Ctrl+D (4) to exit
		if n == 1 && (buf[0] == 3 || buf[0] == 4) {
			return
		}

		// Backspace (127)
		if n == 1 && buf[0] == 127 {
			if col > 0 {
				cell := grid[row][col]
				cell.Status = test.Queued
				col--

				prev := grid[row][col]
				prev.Status = test.Active
				_, _ = fmt.Fprint(out, CSI+"1D")
				PrintCell(out, prev)
				PrintCell(out, cell)
				_, _ = fmt.Fprint(out, CSI+"2D")
			}

			continue
		}

		if done {
			// Escape to quit after a test
			if n == 1 && buf[0] == 27 {
				_, _ = fmt.Fprint(out, CSI+"1A")
				_, _ = fmt.Fprint(out, CSI+"2K")

				return
			}

			// Enter to start the next test
			if n == 1 && buf[0] == 13 {
				row = 0
				col = 0
				done = false

				_, _ = fmt.Fprint(out, CSI+"1A")
				_, _ = fmt.Fprint(out, CSI+"2K---\r\n")

				grid = NextGrid(words)
				PrintGrid(out, grid)
			}

			continue // ignore all other inputs
		}

		// Tab to start a fresh test
		if n == 1 && buf[0] == 9 {
			if row > 0 {
				_, _ = fmt.Fprint(out, CSI+strconv.Itoa(row)+"A")
			}

			_, _ = fmt.Fprint(out, "\r"+CSI+"0J")

			row = 0
			col = 0
			grid = NextGrid(words)
			PrintGrid(out, grid)
		}

		if utf8.FullRune(buf) && buf[0] > 31 {
			r, _ := utf8.DecodeRune(buf)

			if startTime.IsZero() {
				startTime = time.Now()
			}

			cell := grid[row][col]
			if cell == nil {
				continue
			}

			cell.Inputs = append(cell.Inputs, r)
			if r == cell.Rune {
				cell.Status = test.Passed
			} else {
				cell.Status = test.Failed
			}

			PrintCell(out, cell)

			col++
			if col == len(grid[row]) || grid[row][col] == nil {
				if row < len(grid)-1 {
					_, _ = fmt.Fprint(out, "\r\n")
					col = 0
					row++

					continue
				}

				duration := time.Since(startTime)
				result := test.Calc(duration, grid).String()

				done = true
				startTime = time.Time{}

				fmt.Fprint(out, "\r\nResult: "+result+"\r\n")
				fmt.Fprint(out, "\r\n[ENTER] next or [ESC] to quit\r\n")
			}
		}
	}
}

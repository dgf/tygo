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

	"github.com/dgf/tygo/internal/config"
	"github.com/dgf/tygo/internal/dict"
	"github.com/dgf/tygo/internal/gen"
	"github.com/dgf/tygo/internal/test"
	"golang.org/x/term"
)

const (
	CSI               = "\033["
	Reset             = CSI + "0m"
	MoveToStart       = "\r" + CSI + "0J"
	StyleActive       = CSI + "7m"
	StyleFailed       = CSI + "38;5;197m"
	StylePassed       = CSI + "2m"
	EraseLineToEnd    = CSI + "2K"
)

func CursorUp(n int) string {
	return CSI + strconv.Itoa(n) + "A"
}

func CursorBack(n int) string {
	return CSI + strconv.Itoa(n) + "D"
}

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
	_, _ = fmt.Fprint(out, CursorUp(len(grid))+"\r")
}

func NextGrid(cfg config.Config, words []string) test.Grid {
	list := gen.WeightedRandomList(cfg.WordCount, words)

	if cfg.Numbers {
		list = gen.WithNumbers(cfg.Distribution.Number, list)
	}

	if cfg.Punctuation {
		list = gen.PunctuationMarks(list, map[gen.Punctuation]int{
			gen.Word:        cfg.Distribution.Word,
			gen.Period:      cfg.Distribution.Period,
			gen.Comma:       cfg.Distribution.Comma,
			gen.Quotation:   cfg.Distribution.Quotation,
			gen.Question:    cfg.Distribution.Question,
			gen.Exclamation: cfg.Distribution.Exclamation,
			gen.Brackets:    cfg.Distribution.Brackets,
			gen.Braces:      cfg.Distribution.Braces,
			gen.Parenthesis: cfg.Distribution.Parenthesis,
			gen.Colon:       cfg.Distribution.Colon,
			gen.Semicolon:   cfg.Distribution.Semicolon,
		})
	}

	return test.ToGrid(cfg.Width-1, list)
}

func ResetGrid(out io.Writer, cfg config.Config, words []string) test.Grid {
	_, _ = fmt.Fprint(out, MoveToStart)

	grid := NextGrid(cfg, words)
	PrintGrid(out, grid)

	return grid
}

func RetractRune(out io.Writer, grid test.Grid, col, row int) int {
	cell := grid[row][col]
	cell.Status = test.Queued
	col--

	prev := grid[row][col]
	prev.Status = test.Active
	_, _ = fmt.Fprint(out, CursorBack(1))
	PrintCell(out, prev)
	PrintCell(out, cell)
	_, _ = fmt.Fprint(out, CursorBack(2))

	return col
}

func PrintResult(out io.Writer, start time.Time, grid test.Grid) {
	duration := time.Since(start)
	result := test.Calc(duration, grid).String()

	fmt.Fprint(out, "\r\n\r\nResult: "+result+"\r\n")
	fmt.Fprint(out, "\r\n[ENTER] next or [ESC] to quit\r\n")
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

func main() {
	var (
		err error
		cfg config.Config
	)

	cfg, err = config.LoadUserConfig()
	if err != nil {
		cfg = config.Default()
		_ = config.WriteUserConfig(cfg)
	}

	var flags struct {
		dictParam string
		dictTop   int

		fileName   string
		strictMode bool

		wordCount   int
		termCols    int
		numbers     bool
		punctuation bool
	}

	flag.StringVar(&flags.dictParam, "dict", cfg.Dictionary, "dictionary to use, available: german, english")
	flag.IntVar(&flags.dictTop, "top", cfg.TopWords, "top count of words to load from source (dict or file)")

	flag.StringVar(&flags.fileName, "file", "", "vocabulary JSON file with 'words' list")
	flag.BoolVar(&flags.strictMode, "strict", cfg.StrictMode, "enable strict mode, restarts on every error")

	flag.IntVar(&flags.wordCount, "count", cfg.WordCount, "number of words to include in the typing test")
	flag.IntVar(&flags.termCols, "width", cfg.Width, "display width for the typing text")
	flag.BoolVar(&flags.numbers, "nums", cfg.Numbers, "enable number mode")
	flag.BoolVar(&flags.punctuation, "punct", cfg.Punctuation, "enable punctuation marks")

	flag.Parse()

	cfg.Dictionary = flags.dictParam
	cfg.TopWords = flags.dictTop

	cfg.StrictMode = flags.strictMode

	cfg.WordCount = flags.wordCount
	cfg.Width = flags.termCols
	cfg.Numbers = flags.numbers
	cfg.Punctuation = flags.punctuation

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
	if len(flags.fileName) > 0 {
		words = dict.LoadFile(flags.fileName)
	} else {
		words = dict.LoadDict(Dictionary(cfg.Dictionary), cfg.TopWords)
	}

	row := 0
	col := 0
	done := false

	grid := NextGrid(cfg, words)
	PrintGrid(out, grid)

	var start time.Time

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

		if done {
			// Escape to quit after a test
			if n == 1 && buf[0] == 27 {
				_, _ = fmt.Fprint(out, CursorUp(1))
				_, _ = fmt.Fprint(out, EraseLineToEnd)

				return
			}

			// Enter to start the next test
			if n == 1 && buf[0] == 13 {
				row = 0
				col = 0
				done = false

				_, _ = fmt.Fprint(out, CursorUp(1))
				_, _ = fmt.Fprint(out, EraseLineToEnd+"---\r\n")

				grid = NextGrid(cfg, words)
				PrintGrid(out, grid)
			}

			continue // ignore all other inputs
		}

		if !done {
			// Backspace (127) => delete rune
			if n == 1 && buf[0] == 127 {
				if col > 0 {
					col = RetractRune(out, grid, col, row)
				}

				continue
			}

			// Ctrl+W (23) => delete word
			if n == 1 && buf[0] == 23 {
				for col > 0 {
					col = RetractRune(out, grid, col, row)

					if col > 1 && grid[row][col-1].Rune == ' ' {
						break
					}
				}

				continue
			}

			// Tab to start a fresh test
			if n == 1 && buf[0] == 9 {
				if row > 0 {
					_, _ = fmt.Fprint(out, CursorUp(row))
				}

				_, _ = fmt.Fprint(out, MoveToStart)

				grid = ResetGrid(out, cfg, words)

				row = 0
				col = 0
				start = time.Time{}

				continue // ignore all other inputs
			}
		}

		if utf8.FullRune(buf) && buf[0] > 31 {
			r, _ := utf8.DecodeRune(buf)

			if start.IsZero() {
				start = time.Now()
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

				if cfg.StrictMode {
					PrintCell(out, cell)

					row++
					for row < len(grid) {
						row++

						fmt.Fprintf(out, "\r\n")
					}

					PrintResult(out, start, grid)

					done = true
					start = time.Time{}

					continue // ignore all other inputs
				}
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

				PrintResult(out, start, grid)

				done = true
				start = time.Time{}
			}
		}
	}
}

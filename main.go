package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dgf/tygo/internal/dict"
	"github.com/dgf/tygo/internal/gen"
	"golang.org/x/term"
)

type Status int

const (
	Queued Status = iota
	Failed
	Passed
	Active
)

type Cell struct {
	Inputs []rune
	Rune   rune
	Status Status
}

func (c *Cell) String() string {
	return fmt.Sprintf("%d %s %v", c.Status, string(c.Rune), c.Inputs)
}

type Grid [][]*Cell

func (g Grid) String() string {
	sb := strings.Builder{}

	for _, row := range g {
		cells := make([]string, len(row))

		for c, cell := range row {
			if cell != nil {
				cells[c] = cell.String()
			}
		}

		_, _ = sb.WriteString(strings.Join(cells, " "))
		_, _ = sb.WriteString("\r\n")
	}

	return sb.String()
}

type Result struct {
	Duration               time.Duration
	WordsPerMinute         int // WPM = (total keys pressed / 5) / duration in minutes
	AccuracyPercent        int // AP = (correct keys pressed / total keys pressed) * 100
	AdjustedWordsPerMinute int // AWPM = WPM * AP
}

func (r Result) String() string {
	return fmt.Sprintf("%s\r\nWPM %4d\r\nACC  %3d%%\r\nAWPM %3d",
		r.Duration, r.WordsPerMinute, r.AccuracyPercent, r.AdjustedWordsPerMinute)
}

func CalcResults(duration time.Duration, grid Grid) Result {
	totalKeysPressed := 0
	correctKeysPressed := 0

	for _, row := range grid {
		for _, cell := range row {
			if cell == nil {
				break
			}

			totalKeysPressed += len(cell.Inputs)

			for _, i := range cell.Inputs {
				if i == cell.Rune {
					correctKeysPressed++
				}
			}
		}
	}

	wpm := (float64(totalKeysPressed / 5)) / duration.Minutes()
	accuracy := float64(correctKeysPressed) / float64(totalKeysPressed)

	return Result{
		Duration:               duration,
		WordsPerMinute:         int(wpm),
		AccuracyPercent:        int(100 * accuracy),
		AdjustedWordsPerMinute: int(wpm * accuracy),
	}
}

func ToGrid(cols int, wordLines [][][]rune) Grid {
	grid := make(Grid, len(wordLines))

	for l, line := range wordLines {
		grid[l] = make([]*Cell, cols)

		lc := 0

		for w, word := range line {
			for _, r := range word {
				grid[l][lc] = &Cell{Rune: r}
				lc++
			}

			if w < len(line)-1 || l < len(grid)-1 {
				grid[l][lc] = &Cell{Rune: ' '}
				lc++
			}
		}
	}

	return grid
}

func ToLines(cols int, words []string) [][][]rune {
	lines := [][][]rune{}
	line := [][]rune{}
	lc := 0

	for _, word := range words {
		runes := []rune(word)
		if cols < lc+len(runes) {
			lines = append(lines, line)
			line = [][]rune{}
			lc = 0
		}

		line = append(line, runes)
		lc += len(word) + 1
	}

	if lc > 0 {
		lines = append(lines, line)
	}

	return lines
}

const (
	CSI           = "\033["
	Reset         = CSI + "0m"
	RestoreCursor = CSI + "u"
	SaveCursor    = CSI + "s"
	StyleActive   = CSI + "7m"
	StyleFailed   = CSI + "38;5;197m"
	StylePassed   = CSI + "2m"
)

func ColorCSI(s Status) string {
	switch s {
	case Active:
		return StyleActive
	case Failed:
		return StyleFailed
	case Passed:
		return StylePassed
	case Queued:
		return ""
	default:
		return ""
	}
}

func PrintCell(out io.Writer, c *Cell) {
	r := c.Rune

	if r == ' ' && c.Status == Failed {
		r = '_'
	}

	_, _ = fmt.Fprint(out, ColorCSI(c.Status)+string(r)+Reset)
}

func PrintGrid(out io.Writer, grid Grid) {
	for _, row := range grid {
		for c, cell := range row {
			if cell == nil || c == len(row)-1 {
				_, _ = fmt.Fprint(out, "\r\n")

				break
			}

			PrintCell(out, cell)
		}
	}

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

func NextGrid(words []string) Grid {
	test := gen.WeightedRandomList(wordCount, words)

	if numbers {
		test = gen.WithNumbers(test)
	}

	if punctuation {
		test = gen.PunctuationMarks(test)
	}

	lines := ToLines(termCols-1, test)
	grid := ToGrid(termCols, lines)

	return grid
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
				cell.Status = Queued
				col--

				prev := grid[row][col]
				prev.Status = Active
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
				cell.Status = Passed
			} else {
				cell.Status = Failed
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
				result := CalcResults(duration, grid).String()

				done = true
				startTime = time.Time{}

				fmt.Fprint(out, "\r\nResult: "+result+"\r\n")
				fmt.Fprint(out, "\r\n[ENTER] next or [ESC] to quit\r\n")
			}
		}
	}
}

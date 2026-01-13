package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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

type Status int

const (
	Queued Status = iota
	Failed
	Passed
	Active
)

func (s Status) Color() string {
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

type Cell struct {
	Inputs []rune
	Rune   rune
	Status Status
}

func (c *Cell) Render() string {
	r := c.Rune

	if r == ' ' && c.Status == Failed {
		r = '_'
	}

	return c.Status.Color() + string(r) + Reset
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

func WeightedRandom(amount int, words []string) []string {
	result := make([]string, amount)
	count := len(words)
	sum := (count * (count + 1)) / 2

	for a := range amount {
		n := rand.Intn(sum)

		calls := 0
		s := sort.Search(count, func(c int) bool {
			calls++

			return ((c+1)*(c+2))/2 >= n
		})

		result[a] = words[count-s-1]
	}

	return result
}

func LoadFile(fileName string) []string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	type languageFile struct {
		Words []string `json:"words"`
	}

	var lf languageFile

	err = json.Unmarshal(data, &lf)
	if err != nil {
		panic(err)
	}

	return lf.Words
}

var (
	fileName  string
	termCols  int
	wordCount int
)

func init() {
	flag.StringVar(&fileName, "file", "english_1k.json", "vocabulary JSON file with 'words' list")
	flag.IntVar(&wordCount, "count", 20, "number of words to include in the typing test")
	flag.IntVar(&termCols, "width", 50, "display width for the typing text")
	flag.Parse()
}

func main() {
	in := os.Stdin
	inFd := int(in.Fd())

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
		err = term.Restore(inFd, termState)
		if err != nil {
			panic(err)
		}

		os.Exit(0)
	}()

	words := LoadFile(fileName)
	test := WeightedRandom(wordCount, words)
	lines := ToLines(termCols-1, test)
	grid := ToGrid(termCols, lines)

	row := 0
	col := 0
	out := os.Stdout

	for _, row := range grid {
		for c, cell := range row {
			if cell == nil || c == len(row)-1 {
				_, _ = fmt.Fprint(out, "\r\n")

				break
			}

			_, _ = fmt.Fprint(out, cell.Render())
		}
	}

	_, _ = fmt.Fprint(out, CSI+strconv.Itoa(len(grid))+"A\r")

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
				_, _ = fmt.Fprint(out, prev.Render())
				_, _ = fmt.Fprint(out, cell.Render())
				_, _ = fmt.Fprint(out, CSI+"2D")
			}

			continue
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

			_, _ = fmt.Fprint(out, cell.Render())

			col++
			if col == len(grid[row]) || grid[row][col] == nil {
				if row < len(grid)-1 {
					_, _ = fmt.Fprint(out, "\r\n")
					col = 0
					row++

					continue
				}

				duration := time.Since(startTime)
				fmt.Fprint(out, "\r\nResult: "+CalcResults(duration, grid).String()+"\r\n")

				return
			}
		}
	}
}

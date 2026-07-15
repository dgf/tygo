// Tygo is a terminal-based typing test.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"unicode/utf8"

	"github.com/dgf/tygo/internal/config"
	"github.com/dgf/tygo/internal/dict"
	"github.com/dgf/tygo/internal/display"
	"github.com/dgf/tygo/internal/game"
	"github.com/dgf/tygo/internal/gen"
	"github.com/dgf/tygo/internal/test"
	"golang.org/x/term"
)

// Keyboard codes.
const (
	KeyCtrlC       = 3
	KeyCtrlD       = 4
	KeyTab         = 9
	KeyEnter       = 13
	KeyCtrlW       = 23
	KeyEscape      = 27
	KeyBackspace   = 127
	MaxControlCode = 31
)

// Exit codes.
const (
	ExitSuccess          = 0
	ExitUserError        = 1
	ExitEnvironmentError = 2
	ExitInternalError    = 3
)

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

func MustLoadWords(cfg config.Config, file string) []string {
	if len(file) == 0 {
		return dict.LoadDict(Dictionary(cfg.Dictionary), cfg.TopWords)
	}

	words, err := dict.LoadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dictionary load failed: %v\n", err)
		os.Exit(ExitUserError)
	}

	return words
}

func MustMakeRaw(in *os.File) (int, *term.State) {
	fd := int(in.Fd())

	if !term.IsTerminal(fd) {
		fmt.Fprintln(os.Stderr, "Use a terminal (requires a TTY)")
		os.Exit(ExitUserError)
	}

	termState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Raw mode activation failed: %v\n", err)
		os.Exit(ExitEnvironmentError)
	}

	return fd, termState
}

type Actions struct {
	Advance     func(*game.Session, rune)
	Exit        func()
	Next        func(*game.Session)
	Print       func(*game.Session)
	Reset       func(upRows int, session *game.Session)
	RetractRune func(*game.Session)
	RetractWord func(*game.Session)
}

func NewActions(out *os.File) Actions {
	return Actions{
		Advance: func(s *game.Session, r rune) {
			cell, newline := s.Advance(r)

			if cell != nil {
				display.PrintCell(out, cell)
			}

			if newline {
				display.NewLine(out)
			}

			if s.Done() {
				remainRows := len(s.Grid()) - s.Row()
				if remainRows > 1 {
					display.CursorDown(out, remainRows-1)
				}

				result := test.Calc(s.Duration(), s.Grid())
				display.PrintResult(out, result)
			}
		},

		Exit: func() {
			display.UndoLine(out)
		},

		Next: func(s *game.Session) {
			display.UndoLine(out)
			display.PrintLine(out, "---")
			display.PrintGrid(out, s.Grid())
		},

		Print: func(s *game.Session) {
			display.PrintGrid(out, s.Grid())
		},

		Reset: func(upRows int, s *game.Session) {
			if upRows > 0 {
				display.CursorUp(out, upRows)
			}

			display.ResetGrid(out, s.Grid())
		},

		RetractRune: func(s *game.Session) {
			curr, next := s.RetractRune()

			if curr != nil {
				display.CursorBack(out, 1)
				display.PrintCell(out, curr)
				display.PrintCell(out, next)
				display.CursorBack(out, 2)
			}
		},

		RetractWord: func(s *game.Session) {
			cells := s.RetractWord()

			if len(cells) > 0 {
				display.CursorBack(out, len(cells)-1)

				for _, cell := range cells {
					display.PrintCell(out, cell)
				}

				display.CursorBack(out, len(cells))
			}
		},
	}
}

func NewSession(cfg config.Config, words []string) *game.Session {
	list := gen.SampleWeightedList(cfg.WordCount, 5, words)

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

	grid := test.ToGrid(cfg.Width-1, list)

	return game.New(cfg.StrictMode, grid)
}

func RunInputLoop(in *os.File, cfg config.Config, words []string, act Actions) {
	buf := make([]byte, 4)
	session := NewSession(cfg, words)

	act.Print(session)

	for {
		n, err := in.Read(buf)
		if err != nil {
			return // stdin closed > time to leave
		}

		if buf[0] > MaxControlCode && utf8.FullRune(buf[:n]) {
			r, _ := utf8.DecodeRune(buf[:n])

			act.Advance(session, r)
		}

		if buf[0] == KeyCtrlC || buf[0] == KeyCtrlD {
			return
		}

		if session.Done() {
			switch buf[0] {
			case KeyEscape:
				act.Exit()

				return
			case KeyEnter:
				session = NewSession(cfg, words)

				act.Next(session)
			}

			continue // ignore all other inputs
		}

		switch buf[0] {
		case KeyBackspace:
			act.RetractRune(session)
		case KeyCtrlW:
			act.RetractWord(session)
		case KeyTab:
			upRows := session.Row()
			session = NewSession(cfg, words)

			act.Reset(upRows, session)
		}
	}
}

func main() {
	cfg, loadErr := config.LoadUserConfig()
	if loadErr != nil {
		cfg = config.Default()
		_ = config.WriteUserConfig(cfg)
	}

	var file string

	flag.StringVar(&cfg.Dictionary, "dict", cfg.Dictionary, "dictionary to use, available: german, english")

	flag.IntVar(&cfg.TopWords, "top", cfg.TopWords, "top count of words to load from source (dict or file)")
	flag.IntVar(&cfg.WordCount, "count", cfg.WordCount, "number of words to include in the typing test")
	flag.IntVar(&cfg.Width, "width", cfg.Width, "display width for the typing text")

	flag.BoolVar(&cfg.Numbers, "nums", cfg.Numbers, "enable number mode")
	flag.BoolVar(&cfg.Punctuation, "punct", cfg.Punctuation, "enable punctuation marks")
	flag.BoolVar(&cfg.StrictMode, "strict", cfg.StrictMode, "enable strict mode, restarts on every error")

	flag.StringVar(&file, "file", "", "vocabulary JSON file with 'words' list")

	flag.Parse()

	in := os.Stdin
	out := os.Stdout
	words := MustLoadWords(cfg, file)
	fd, oldState := MustMakeRaw(in)

	defer func() {
		_ = term.Restore(fd, oldState)

		if r := recover(); r != nil {
			fmt.Fprintf(out, "%v - %s", r, debug.Stack())
			os.Exit(ExitInternalError)
		}
	}()

	RunInputLoop(in, cfg, words, NewActions(out))
}

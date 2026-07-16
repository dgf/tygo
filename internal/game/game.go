// Package game manages a typing session and drives the input dispatch loop.
package game

import (
	"io"
	"unicode/utf8"

	"github.com/dgf/tygo/internal/config"
	"github.com/dgf/tygo/internal/gen"
	"github.com/dgf/tygo/internal/test"
)

func newGameSession(cfg config.Config, words []string) *Session {
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

	return NewSession(cfg.StrictMode, grid)
}

func Run(in io.Reader, cfg config.Config, words []string, act Actions) {
	buf := make([]byte, 4)
	session := newGameSession(cfg, words)

	act.Print(session)

	for {
		n, err := in.Read(buf)
		if err != nil {
			return // stdin closed > time to leave
		}

		if n == 1 {
			if buf[0] == KeyCtrlC || buf[0] == KeyCtrlD {
				return
			}

			if session.Done() {
				switch buf[0] {
				case KeyEscape:
					act.Exit(session)

					return
				case KeyEnter:
					session = newGameSession(cfg, words)

					act.Next(session)
				}

				continue
			}

			switch buf[0] {
			case KeyBackspace:
				act.RetractRune(session)

				continue
			case KeyCtrlW:
				act.RetractWord(session)

				continue
			case KeyTab:
				old := session
				session = newGameSession(cfg, words)
				act.Reset(old, session)

				continue
			}
		}

		if buf[0] > MaxControlCode && utf8.FullRune(buf[:n]) {
			r, _ := utf8.DecodeRune(buf[:n])

			act.Advance(session, r)
		}
	}
}

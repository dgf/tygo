// Package game manages a typing session and drives the input dispatch loop.
package game

import (
	"github.com/dgf/tygo/internal/config"
	"github.com/dgf/tygo/internal/gen"
	"github.com/dgf/tygo/internal/test"
)

type Game struct {
	factory  SessionFactory
	renderer Renderer
	session  *Session
}

func NewGame(cfg config.Config, words []string, renderer Renderer) *Game {
	factory := func() *Session {
		return newGameSession(cfg, words)
	}

	session := factory()
	renderer.Reset(session.Grid())

	return &Game{
		factory:  factory,
		renderer: renderer,
		session:  session,
	}
}

func EventActions() map[test.Event]Action {
	return map[test.Event]Action{
		test.EventBackRune: BackRune,
		test.EventBackWord: BackWord,
		test.EventExit:     Exit,
		test.EventNext:     Next,
		test.EventQuit:     Quit,
		test.EventReset:    Reset,
	}
}

func (g *Game) HandleEvent(e test.Event) bool {
	action, ok := EventActions()[e]

	if !ok {
		return false
	}

	g.session = action(g.session, g.renderer, g.factory)

	return g.session == nil
}

func (g *Game) HandleRune(r rune) {
	if g.session.Done() {
		return
	}

	cell, br := g.session.Advance(r)
	g.renderer.Advance(cell, br)

	if g.session.Done() {
		result := test.Calc(g.session.Duration(), g.session.Grid())

		g.renderer.Print(result)
	}
}

func newGameSession(cfg config.Config, words []string) *Session {
	list := gen.SampleWeightedList(cfg.WordCount, cfg.NoRepeat, words)

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

	return NewSession(cfg.StrictMode, test.ToGrid(cfg.Width-1, list))
}

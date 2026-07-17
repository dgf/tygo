package game

import "github.com/dgf/tygo/internal/test"

type Action func(s *Session, r Renderer, f SessionFactory) *Session

func BackRune(s *Session, r Renderer, _ SessionFactory) *Session {
	if !s.Done() {
		curr, next := s.RetractRune()

		if curr != nil {
			r.Retract(test.Cells{curr, next})
		}
	}

	return s
}

func BackWord(s *Session, r Renderer, _ SessionFactory) *Session {
	if !s.Done() {
		cells := s.RetractWord()

		if len(cells) > 0 {
			r.Retract(cells)
		}
	}

	return s
}

func Exit(_ *Session, _ Renderer, _ SessionFactory) *Session {
	return nil
}

func Next(s *Session, r Renderer, f SessionFactory) *Session {
	if !s.Done() {
		return s
	}

	n := f()
	r.Next(n.Grid())

	return n
}

func Quit(s *Session, r Renderer, _ SessionFactory) *Session {
	if !s.Done() {
		return s
	}

	r.Exit()

	return nil
}

func Reset(s *Session, r Renderer, f SessionFactory) *Session {
	if s.Done() {
		return s
	}

	n := f()
	r.Reset(n.Grid())

	return n
}

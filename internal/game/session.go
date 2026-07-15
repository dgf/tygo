package game

import (
	"slices"
	"time"

	"github.com/dgf/tygo/internal/test"
)

type Session struct {
	strict bool
	row    int
	col    int

	duration time.Duration
	start    time.Time
	grid     test.Grid
}

func NewSession(strict bool, grid test.Grid) *Session {
	return &Session{
		grid:     grid,
		strict:   strict,
		start:    time.Time{},
		duration: 0,
		row:      0,
		col:      0,
	}
}

func (s *Session) Done() bool {
	return s.duration > 0
}

func (s *Session) Duration() time.Duration {
	return s.duration
}

func (s *Session) Grid() test.Grid {
	return s.grid
}

func (s *Session) Row() int {
	return s.row
}

func (s *Session) Advance(r rune) (*test.Cell, bool) {
	if s.Done() {
		return nil, false
	}

	if s.start.IsZero() {
		s.start = time.Now()
	}

	cell := s.grid[s.row][s.col]
	if cell == nil {
		return nil, true
	}

	cell.Inputs = append(cell.Inputs, r)

	if r == cell.Rune {
		cell.Status = test.Passed
	} else {
		cell.Status = test.Failed

		if s.strict {
			s.duration = time.Since(s.start)

			return cell, false
		}
	}

	s.col++
	if s.col == len(s.grid[s.row]) || s.grid[s.row][s.col] == nil {
		if s.row == len(s.grid)-1 {
			s.duration = time.Since(s.start)

			return cell, false
		}

		s.col = 0
		s.row++

		return cell, true
	}

	return cell, false
}

func (s *Session) RetractRune() (*test.Cell, *test.Cell) {
	if s.col < 1 {
		return nil, nil
	}

	curr := s.grid[s.row][s.col]
	curr.Status = test.Queued
	s.col--

	prev := s.grid[s.row][s.col]
	prev.Status = test.Active

	return prev, curr
}

func (s *Session) RetractWord() []*test.Cell {
	if s.col < 1 {
		return []*test.Cell{}
	}

	currs := []*test.Cell{}

	for s.col > 0 {
		curr := s.grid[s.row][s.col]
		curr.Status = test.Queued
		currs = append(currs, curr)

		s.col--

		if s.col > 1 && s.grid[s.row][s.col-1].Rune == ' ' {
			break
		}
	}

	curr := s.grid[s.row][s.col]
	curr.Status = test.Active
	currs = append(currs, curr)

	slices.Reverse(currs)

	return currs
}

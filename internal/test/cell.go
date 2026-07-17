package test

import "fmt"

type Cell struct {
	Inputs []rune
	Rune   rune
	Status Status
}

type Cells []*Cell

func (c Cell) String() string {
	return fmt.Sprintf("rune: %q, status: %v", c.Rune, c.Status)
}

func Enqueue(r rune) *Cell {
	return &Cell{Rune: r, Status: Queued, Inputs: []rune{}}
}

package test

type Cell struct {
	Inputs []rune
	Rune   rune
	Status Status
}

func Enqueue(r rune) *Cell {
	return &Cell{Rune: r, Status: Queued, Inputs: []rune{}}
}

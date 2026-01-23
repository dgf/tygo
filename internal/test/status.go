package test

type Status int

const (
	Queued Status = iota
	Failed
	Passed
	Active
)

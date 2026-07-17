package test

type Event int

// Events for input dispatch.
const (
	EventBackRune Event = iota
	EventBackWord
	EventExit
	EventNext
	EventQuit
	EventReset
)

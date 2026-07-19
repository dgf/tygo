// Package input reads keyboard events and dispatches them to a Handler.
package input

import "github.com/dgf/tygo/internal/test"

type Handler interface {
	HandleEvent(e test.Event) (quit bool)
	HandleRune(r rune)
}

func KeyEvents() map[KeyCode]test.Event {
	return map[KeyCode]test.Event{
		KeyCtrlC:     test.EventExit,
		KeyCtrlD:     test.EventExit,
		KeyEscape:    test.EventQuit,
		KeyEnter:     test.EventNext,
		KeyBackspace: test.EventBackRune,
		KeyCtrlW:     test.EventBackWord,
		KeyTab:       test.EventReset,
	}
}

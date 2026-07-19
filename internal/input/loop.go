package input

import (
	"io"
	"unicode/utf8"

	"github.com/dgf/tygo/internal/test"
)

func Loop(in io.Reader, handler Handler) {
	buf := make([]byte, 4)
	quit := handler.HandleEvent(test.EventNext)

	for !quit {
		n, err := in.Read(buf)
		if err != nil {
			return // stdin closed > time to leave
		}

		if n == 1 {
			e, ok := KeyEvents()[KeyCode(buf[0])]

			if ok {
				quit = handler.HandleEvent(e)

				continue
			}
		}

		if buf[0] > MaxControlCode && utf8.FullRune(buf[:n]) {
			r, _ := utf8.DecodeRune(buf[:n])

			handler.HandleRune(r)
		}
	}
}

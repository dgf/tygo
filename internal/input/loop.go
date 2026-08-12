package input

import (
	"io"
	"unicode/utf8"

	"github.com/dgf/tygo/internal/test"
)

const InputBufferSize = 4

func Loop(in io.Reader, handler Handler) {
	buf := make([]byte, InputBufferSize)
	quit := handler.HandleEvent(test.EventNext)

	for !quit {
		count, err := in.Read(buf)
		if err != nil {
			return // stdin closed > time to leave
		}

		if count == 1 {
			e, ok := KeyEvents()[KeyCode(buf[0])]

			if ok {
				quit = handler.HandleEvent(e)

				continue
			}
		}

		if buf[0] > MaxControlCode && utf8.FullRune(buf[:count]) {
			r, _ := utf8.DecodeRune(buf[:count])

			handler.HandleRune(r)
		}
	}
}

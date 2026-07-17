package input

type KeyCode int

// Keyboard codes.
const (
	KeyCtrlC     KeyCode = 3
	KeyCtrlD     KeyCode = 4
	KeyTab       KeyCode = 9
	KeyEnter     KeyCode = 13
	KeyCtrlW     KeyCode = 23
	KeyEscape    KeyCode = 27
	KeyBackspace KeyCode = 127

	MaxControlCode = 31
)

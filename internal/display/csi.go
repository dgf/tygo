package display

// CSI sequences.
const (
	CSI            = "\033["
	Reset          = CSI + "0m"
	EraseLineToEnd = CSI + "2K"
	MoveToStart    = CSI + "0J"
)

// Text styles.
const (
	StylePassed = CSI + "2m"
	StyleActive = CSI + "7m"
	StyleFailed = CSI + "38;5;197m"
)

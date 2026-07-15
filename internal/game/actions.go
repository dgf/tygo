package game

// Actions bundles the input callbacks.
type Actions struct {
	Advance     func(*Session, rune)
	Exit        func()
	Next        func(*Session)
	Print       func(*Session)
	Reset       func(upRows int, state *Session)
	RetractRune func(*Session)
	RetractWord func(*Session)
}

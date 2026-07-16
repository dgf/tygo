package game

type (
	Advance func(*Session, rune)
	Action  func(*Session)
	Reset   func(old, session *Session)
)

type Actions struct {
	Advance     Advance
	Exit        Action
	Next        Action
	Print       Action
	Reset       Reset
	RetractRune Action
	RetractWord Action
}

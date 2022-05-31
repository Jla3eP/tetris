package constants

const (
	ColorBlue = iota
	ColorGreen
	ColorOrange
	ColorRed
	ColorYellow
	ColorBlack
)

const (
	StatusProcessing = iota
	StatusWaiting
	StatusPlaying
	StatusWatching
	StatusEnd
)

type Event int

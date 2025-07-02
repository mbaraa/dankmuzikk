package models

type PlayerLoopMode string

const (
	LoopAllMode  PlayerLoopMode = "all"
	LoopOnceMode PlayerLoopMode = "once"
	LoopOffMode  PlayerLoopMode = "off"
)

type PlayerState struct {
	Shuffled          bool
	CurrentSongIndex  int
	LoopMode          PlayerLoopMode
	CurrentPlaylistId uint
	Songs             []Song
}

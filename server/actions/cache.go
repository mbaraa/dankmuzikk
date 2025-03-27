package actions

type Cache interface {
	StoreLyrics(songId uint, lyrics []string) error
	GetLyrics(songId uint) ([]string, error)
}

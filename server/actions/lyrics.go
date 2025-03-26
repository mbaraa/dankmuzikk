package actions

type Lyrics interface {
	GetForSong(songName string) ([]string, error)
	GetForSongAndArtist(songName, artistName string) ([]string, error)
}

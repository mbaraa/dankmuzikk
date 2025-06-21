package actions

type Lyrics interface {
	GetForSong(songName string) ([]string, map[string]string, error)
	GetForSongAndArtist(songName, artistName string) ([]string, map[string]string, error)
}

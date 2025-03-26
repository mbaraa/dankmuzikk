package genius

import (
	"dankmuzikk/app"
	"dankmuzikk/config"
	"fmt"

	"github.com/mbaraa/gonius"
)

type GeniusLyrics struct {
	client *gonius.Client
}

func New() *GeniusLyrics {
	return &GeniusLyrics{
		client: gonius.NewClient(config.Env().GeniusToken),
	}
}

func (l *GeniusLyrics) GetForSong(songName string) ([]string, error) {
	return l.getSongLyrics(searchInput{
		SongName: songName,
	})
}

func (l *GeniusLyrics) GetForSongAndArtist(songName, artistName string) ([]string, error) {
	return l.getSongLyrics(searchInput{
		SongName:   songName,
		ArtistName: artistName,
	})
}

type searchInput struct {
	SongName   string
	ArtistName string
	AlbumName  string
}

func (l *GeniusLyrics) getSongLyrics(s searchInput) ([]string, error) {
	var hits []gonius.Hit
	var err error

	okArtist := s.ArtistName != ""
	okAlbum := s.AlbumName != ""
	okSong := s.SongName != ""

	switch {
	case !okArtist && !okAlbum && okSong:
		hits, err = l.client.Search.Get(s.SongName)
		if err != nil {
			return nil, err
		}
	case okArtist && !okAlbum && okSong:
		hits, err = l.client.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.ArtistName))
		if err != nil {
			return nil, err
		}
	case !okArtist && okAlbum && okSong:
		hits, err = l.client.Search.Get(fmt.Sprintf("%s %s", s.SongName, s.AlbumName))
		if err != nil {
			return nil, err
		}
	case okArtist && okAlbum && okSong:
		hits, err = l.client.Search.Get(fmt.Sprintf("%s %s %s", s.SongName, s.AlbumName, s.ArtistName))
		if err != nil {
			return nil, err
		}
	}

	if len(hits) == 0 {
		return nil, &app.ErrNotFound{
			ResourceName: "lyrics",
		}
	}

	lyrics, err := l.client.Lyrics.FindForSong(hits[0].Result.URL)
	if err != nil {
		return nil, err
	}

	return lyrics.Parts(), nil
}

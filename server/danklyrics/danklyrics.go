package danklyrics

import (
	"dankmuzikk/app"
	"dankmuzikk/config"

	"github.com/mbaraa/danklyrics/pkg/client"
	"github.com/mbaraa/danklyrics/pkg/provider"
)

type dankLyrics struct {
	client *client.Http
}

func New() *dankLyrics {
	client, _ := client.NewHttp(client.Config{
		Providers:  []provider.Name{provider.Dank, provider.LyricFind},
		ApiAddress: config.Env().DankLyricsAddress,
	})

	return &dankLyrics{
		client: client,
	}
}

func (d *dankLyrics) GetForSong(songName string) ([]string, map[string]string, error) {
	lyrics, err := d.client.GetSongLyrics(provider.SearchParams{
		SongName: songName,
	})
	if err != nil {
		return nil, nil, err
	}

	if len(lyrics.Parts) == 0 && len(lyrics.Synced) == 0 {
		return nil, nil, &app.ErrNotFound{
			ResourceName: "lyrics",
		}
	}

	return lyrics.Parts, lyrics.Synced, nil
}

func (d *dankLyrics) GetForSongAndArtist(songName, artistName string) ([]string, map[string]string, error) {
	lyrics, err := d.client.GetSongLyrics(provider.SearchParams{
		SongName:   songName,
		ArtistName: artistName,
	})
	if err != nil {
		return nil, nil, err
	}

	if len(lyrics.Parts) == 0 && len(lyrics.Synced) == 0 {
		return nil, nil, &app.ErrNotFound{
			ResourceName: "lyrics",
		}
	}

	return lyrics.Parts, lyrics.Synced, nil
}

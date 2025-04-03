package actions

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"dankmuzikk/log"
	"errors"
	"fmt"
	"regexp"
)

type Song struct {
	PublicId        string `json:"public_id"`
	Title           string `json:"title"`
	Artist          string `json:"artist"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	Duration        string `json:"duration"`
	PlayTimes       int    `json:"play_times,omitempty"`
	Votes           int    `json:"votes,omitempty"`
	AddedAt         string `json:"added_at,omitempty"`
	FullyDownloaded bool   `json:"fully_downloaded"`
}

type GetSongByPublicIdParams struct {
	SongPublicId string
}

func (a *Actions) GetSongByPublicId(params GetSongByPublicIdParams) (Song, error) {
	song, err := a.app.GetSongByPublicId(params.SongPublicId)
	if err != nil && errors.As(err, &app.ErrNotFound{}) {
		searches, err := a.SearchYouTube(params.SongPublicId)
		if err != nil {
			return Song{}, err
		}
		if len(searches) == 0 {
			return Song{}, &app.ErrNotFound{
				ResourceName: "song",
			}
		}
		for i, s := range searches {
			if s.PublicId == params.SongPublicId {
				ss := models.Song{
					PublicId:        s.PublicId,
					Title:           s.Title,
					Artist:          s.Artist,
					ThumbnailUrl:    fmt.Sprintf("%s/pix/%s.webp", config.CdnAddress(), s.PublicId),
					Duration:        s.Duration,
					FullyDownloaded: false,
				}
				newSong, err := a.app.CreateSong(ss)
				if err != nil {
					return Song{}, err
				}

				if i == 0 {
					song = newSong
				}
			}
		}
		err = a.eventhub.Publish(events.SongPlayed{
			SongPublicId: params.SongPublicId,
		})
		if err != nil {
			return Song{}, err
		}
	} else if err != nil {
		return Song{}, err
	}

	return Song{
		PublicId:        song.PublicId,
		Title:           song.Title,
		Artist:          song.Artist,
		ThumbnailUrl:    song.ThumbnailUrl,
		Duration:        song.Duration,
		FullyDownloaded: song.FullyDownloaded,
	}, nil
}

func (a *Actions) IncrementSongPlaysInPlaylist(songId, playlistId string, profileId uint) error {
	return a.app.IncrementSongPlaysInPlaylist(songId, playlistId, profileId)
}

type UpvoteSongInPlaylistParams struct {
	ActionContext    `json:"-"`
	SongPublicId     string
	PlaylistPublicId string
}

type UpvoteSongInPlaylistPayload struct {
	VotesCount int `json:"votes_count"`
}

func (a *Actions) UpvoteSongInPlaylist(params UpvoteSongInPlaylistParams) (UpvoteSongInPlaylistPayload, error) {
	newCount, err := a.app.UpvoteSongInPlaylist(params.SongPublicId, params.PlaylistPublicId, params.Account.Id)
	if err != nil {
		return UpvoteSongInPlaylistPayload{}, err
	}

	return UpvoteSongInPlaylistPayload{
		VotesCount: newCount,
	}, nil
}

type DownvoteSongInPlaylistParams struct {
	ActionContext    `json:"-"`
	SongPublicId     string
	PlaylistPublicId string
}

type DownvoteSongInPlaylistPayload struct {
	VotesCount int `json:"votes_count"`
}

func (a *Actions) DownvoteSongInPlaylist(params DownvoteSongInPlaylistParams) (DownvoteSongInPlaylistPayload, error) {
	newCount, err := a.app.DownvoteSongInPlaylist(params.SongPublicId, params.PlaylistPublicId, params.Account.Id)
	if err != nil {
		return DownvoteSongInPlaylistPayload{}, err
	}

	return DownvoteSongInPlaylistPayload{
		VotesCount: newCount,
	}, nil
}

func (a *Actions) AddSongToHistory(songPublicId string, profileId uint) error {
	if profileId == 0 {
		return nil
	}
	return a.app.AddSongToHistory(songPublicId, profileId)
}

func (a *Actions) DownloadYouTubeSong(ytId string) error {
	song, err := a.app.GetSongByPublicId(ytId)
	if err != nil {
		return err
	}

	if song.FullyDownloaded {
		log.Infof("The song with id %s was already downloaded 😬\n", ytId)
		return nil
	}

	err = a.youtube.DownloadYoutubeSong(ytId)
	if err != nil {
		return err
	}

	return a.eventhub.Publish(events.SongDownloaded{
		SongPublicId: ytId,
	})
}

func (a *Actions) MarkSongAsDownloaded(songPublicId string) error {
	return a.app.MarkSongAsDownloaded(songPublicId)
}

type ToggleSongInPlaylistParams struct {
	ActionContext    `json:"-"`
	SongPublicId     string
	PlaylistPublicId string
}

type ToggleSongInPlaylistPayload struct {
	Added bool `json:"added"`
}

func (a *Actions) ToggleSongInPlaylist(params ToggleSongInPlaylistParams) (ToggleSongInPlaylistPayload, error) {
	added, err := a.app.ToggleSongInPlaylist(params.SongPublicId, params.PlaylistPublicId, params.Account.Id)
	if err != nil {
		return ToggleSongInPlaylistPayload{}, err
	}

	var event events.Event
	if added {
		event = events.SongAddedToPlaylist{
			AccountId:     params.Account.Id,
			PlaylistPubId: params.PlaylistPublicId,
			SongPublicId:  params.SongPublicId,
		}
	} else {
		event = events.SongRemovedFromPlaylist{
			AccountId:     params.Account.Id,
			PlaylistPubId: params.PlaylistPublicId,
			SongPublicId:  params.SongPublicId,
		}
	}
	err = a.eventhub.Publish(event)
	if err != nil {
		return ToggleSongInPlaylistPayload{}, err
	}

	return ToggleSongInPlaylistPayload{
		Added: added,
	}, nil
}

type PlaySongParams struct {
	ActionContext `json:"-"`
	SongPublicId  string
	PlaylistPubId string
}

type PlaySongPayload struct {
	MediaUrl string `json:"media_url"`
}

func (a *Actions) PlaySong(params PlaySongParams) (PlaySongPayload, error) {
	_, err := a.GetSongByPublicId(GetSongByPublicIdParams{
		SongPublicId: params.SongPublicId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:        params.Account.Id,
		SongPublicId:     params.SongPublicId,
		PlaylistPublicId: params.PlaylistPubId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		MediaUrl: fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().CdnAddress, params.SongPublicId),
	}, nil
}

func (a *Actions) SaveSongsMetadataFromYouTube(songs []Song) error {
	for _, newSong := range songs {
		_, _ = a.app.CreateSong(models.Song{
			PublicId:        newSong.PublicId,
			Title:           newSong.Title,
			Artist:          newSong.Artist,
			ThumbnailUrl:    fmt.Sprintf("%s/pix/%s.webp", config.CdnAddress(), newSong.PublicId),
			Duration:        newSong.Duration,
			FullyDownloaded: false,
		})
	}

	return nil
}

type GetLyricsForSongParams struct {
	SongPublicId string
}

type GetLyricsForSongPayload struct {
	SongTitle string   `json:"song_title"`
	Lyrics    []string `json:"lyrics"`
}

var songTitleWeirdStuff = regexp.MustCompile(`(\(.*\)|\[.*\]|\{.*\}|\<.*\>)`)

func (a *Actions) GetLyricsForSong(params GetLyricsForSongParams) (GetLyricsForSongPayload, error) {
	song, err := a.app.GetSongByPublicId(params.SongPublicId)
	if err != nil {
		return GetLyricsForSongPayload{}, err
	}

	lyrics, err := a.cache.GetLyrics(song.Id)
	if lyrics != nil && err == nil {
		return GetLyricsForSongPayload{
			SongTitle: song.Title,
			Lyrics:    lyrics,
		}, nil
	}

	lyrics, err = a.lyrics.GetForSong(songTitleWeirdStuff.ReplaceAllString(song.Title, ""))
	if err != nil {
		return GetLyricsForSongPayload{}, err
	}

	err = a.cache.StoreLyrics(song.Id, lyrics)
	if err != nil {
		return GetLyricsForSongPayload{}, err
	}

	return GetLyricsForSongPayload{
		SongTitle: song.Title,
		Lyrics:    lyrics,
	}, nil
}

func (a *Actions) GetLyricsForSongAndArtist(params GetLyricsForSongParams) (GetLyricsForSongPayload, error) {
	song, err := a.app.GetSongByPublicId(params.SongPublicId)
	if err != nil {
		return GetLyricsForSongPayload{}, err
	}

	lyrics, err := a.cache.GetLyrics(song.Id)
	if lyrics != nil && err == nil {
		return GetLyricsForSongPayload{
			SongTitle: song.Title,
			Lyrics:    lyrics,
		}, nil
	}

	lyrics, err = a.lyrics.GetForSongAndArtist(songTitleWeirdStuff.ReplaceAllString(song.Title, ""), song.Artist)
	if err != nil {
		return GetLyricsForSongPayload{}, err
	}

	err = a.cache.StoreLyrics(song.Id, lyrics)
	if err != nil {
		return GetLyricsForSongPayload{}, err
	}

	return GetLyricsForSongPayload{
		SongTitle: song.Title,
		Lyrics:    lyrics,
	}, nil
}

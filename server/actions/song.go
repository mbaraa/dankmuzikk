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
	YtId            string `json:"yt_id"`
	Title           string `json:"title"`
	Artist          string `json:"artist"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	Duration        string `json:"duration"`
	PlayTimes       int    `json:"play_times,omitempty"`
	Votes           int    `json:"votes,omitempty"`
	AddedAt         string `json:"added_at,omitempty"`
	FullyDownloaded bool   `json:"fully_downloaded"`
}

type GetSongByYouTubeIdParams struct {
	SongYouTubeId string
}

func (a *Actions) GetSongByYouTubeId(params GetSongByYouTubeIdParams) (Song, error) {
	song, err := a.app.GetSongByYouTubeId(params.SongYouTubeId)
	if err != nil && errors.As(err, &app.ErrNotFound{}) {
		searches, err := a.SearchYouTube(params.SongYouTubeId)
		if err != nil {
			return Song{}, err
		}
		if len(searches) == 0 {
			return Song{}, &app.ErrNotFound{
				ResourceName: "song",
			}
		}
		for i, s := range searches {
			if s.YtId == params.SongYouTubeId {
				ss := models.Song{
					YtId:            s.YtId,
					Title:           s.Title,
					Artist:          s.Artist,
					ThumbnailUrl:    fmt.Sprintf("%s/pix/%s.webp", config.CdnAddress(), s.YtId),
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
			SongYtId: params.SongYouTubeId,
		})
		if err != nil {
			return Song{}, err
		}
	} else if err != nil {
		return Song{}, err
	}

	return Song{
		YtId:            song.YtId,
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

func (a *Actions) AddSongToHistory(songYtId string, profileId uint) error {
	return a.app.AddSongToHistory(songYtId, profileId)
}

func (a *Actions) DownloadYouTubeSong(ytId string) error {
	song, err := a.app.GetSongByYouTubeId(ytId)
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
		SongYtId: ytId,
	})
}

func (a *Actions) MarkSongAsDownloaded(songYtId string) error {
	return a.app.MarkSongAsDownloaded(songYtId)
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
			SongYtId:      params.SongPublicId,
		}
	} else {
		event = events.SongRemovedFromPlaylist{
			AccountId:     params.Account.Id,
			PlaylistPubId: params.PlaylistPublicId,
			SongYtId:      params.SongPublicId,
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
	SongYtId      string
	PlaylistPubId string
}

type PlaySongPayload struct {
	MediaUrl string `json:"media_url"`
}

func (a *Actions) PlaySong(params PlaySongParams) (PlaySongPayload, error) {
	_, err := a.GetSongByYouTubeId(GetSongByYouTubeIdParams{
		SongYouTubeId: params.SongYtId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:     params.Account.Id,
		SongYtId:      params.SongYtId,
		PlaylistPubId: params.PlaylistPubId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		MediaUrl: fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().CdnAddress, params.SongYtId),
	}, nil
}

func (a *Actions) SaveSongsMetadataFromYouTube(songs []Song) error {
	for _, newSong := range songs {
		_, _ = a.app.CreateSong(models.Song{
			YtId:            newSong.YtId,
			Title:           newSong.Title,
			Artist:          newSong.Artist,
			ThumbnailUrl:    fmt.Sprintf("%s/pix/%s.webp", config.CdnAddress(), newSong.YtId),
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
	song, err := a.app.GetSongByYouTubeId(params.SongPublicId)
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
	song, err := a.app.GetSongByYouTubeId(params.SongPublicId)
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

package actions

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"dankmuzikk/log"
	"errors"
	"fmt"
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

func (a *Actions) GetSongByYouTubeId(ytId string) (Song, error) {
	song, err := a.app.GetSongByYouTubeId(ytId)
	if err != nil && errors.As(err, &app.ErrNotFound{}) {
		searches, err := a.SearchYouTube(ytId)
		if err != nil {
			return Song{}, err
		}
		if len(searches) == 0 {
			return Song{}, &app.ErrNotFound{
				ResourceName: "song",
			}
		}
		for i, s := range searches {
			if s.YtId == ytId {
				ss := models.Song{
					YtId:            s.YtId,
					Title:           s.Title,
					Artist:          s.Artist,
					ThumbnailUrl:    s.ThumbnailUrl,
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
			SongYtId: ytId,
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

func (a *Actions) UpvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error) {
	return a.app.UpvoteSongInPlaylist(songId, playlistPubId, ownerId)
}

func (a *Actions) DownvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error) {
	return a.app.DownvoteSongInPlaylist(songId, playlistPubId, ownerId)
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

func (a *Actions) ToggleSongInPlaylist(songId, playlistPubId string, ownerId uint) (added bool, err error) {
	added, err = a.app.ToggleSongInPlaylist(songId, playlistPubId, ownerId)
	if err != nil {
		return false, err
	}

	var event events.Event
	if added {
		event = events.SongAddedToPlaylist{
			PlaylistPubId: playlistPubId,
			SongYtId:      songId,
		}
	} else {
		event = events.SongRemovedFromPlaylist{
			PlaylistPubId: playlistPubId,
			SongYtId:      songId,
		}
	}
	err = a.eventhub.Publish(event)
	if err != nil {
		return false, err
	}

	return added, nil
}

type PlaySongParams struct {
	Profile       models.Profile
	SongYtId      string
	PlaylistPubId string
}

func (a *Actions) PlaySong(params PlaySongParams) (mediaUrl string, err error) {
	_, err = a.GetSongByYouTubeId(params.SongYtId)
	if err != nil {
		return "", err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		ProfileId:     params.Profile.Id,
		SongYtId:      params.SongYtId,
		PlaylistPubId: params.PlaylistPubId,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().CdnAddress, params.SongYtId), nil
}

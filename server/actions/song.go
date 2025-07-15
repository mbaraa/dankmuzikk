package actions

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"dankmuzikk/log"
	"errors"
	"fmt"
	"time"
)

type Song struct {
	PublicId     string        `json:"public_id"`
	Title        string        `json:"title"`
	Artist       string        `json:"artist"`
	ThumbnailUrl string        `json:"thumbnail_url"`
	Duration     time.Duration `json:"duration"`
	PlayTimes    int           `json:"play_times,omitempty"`
	Votes        int           `json:"votes,omitempty"`
	AddedAt      string        `json:"added_at,omitempty"`
	MediaUrl     string        `json:"media_url"`
	Favorite     bool          `json:"favorite"`
}

func mapModelToActionsSong(s models.Song) Song {
	mediaUrl := ""
	if s.FullyDownloaded {
		mediaUrl = fmt.Sprintf("%s/muzikkx/%s.mp3", config.CdnAddress(), s.PublicId)
	}
	return Song{
		PublicId:     s.PublicId,
		Title:        s.Title,
		Artist:       s.Artist,
		ThumbnailUrl: s.ThumbnailUrl,
		Duration:     s.RealDuration,
		PlayTimes:    s.PlayTimes,
		Votes:        s.Votes,
		AddedAt:      s.AddedAt,
		Favorite:     s.Favorite,
		MediaUrl:     mediaUrl,
	}
}

type GetSongByPublicIdParams struct {
	ActionContext
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
					RealDuration:    s.Duration,
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
		event := events.SongPlayed{
			SongPublicId: params.SongPublicId,
		}
		err = a.eventhub.Publish(event)
		if err != nil {
			return Song{}, err
		}
	} else if err != nil {
		return Song{}, err
	}

	song.Favorite = a.app.IsSongFavorite(params.Account.Id, song.Id)

	return mapModelToActionsSong(song), nil
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

func (a *Actions) AddSongToHistory(songPublicId string, accountId uint) error {
	if accountId == 0 {
		return nil
	}
	return a.app.AddSongToHistory(songPublicId, accountId)
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
	ActionContext    `json:"-"`
	SongPublicId     string
	PlaylistPublicId string
}

type PlaySongPayload struct {
	Song
}

func (a *Actions) PlaySong(params PlaySongParams) (PlaySongPayload, error) {
	song, err := a.GetSongByPublicId(GetSongByPublicIdParams{
		SongPublicId:  params.SongPublicId,
		ActionContext: params.ActionContext,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:        params.Account.Id,
		ClientHash:       params.ClientHash,
		SongPublicId:     params.SongPublicId,
		PlaylistPublicId: params.PlaylistPublicId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.app.CreateSongsQueue(params.Account.Id, params.ClientHash, []string{params.SongPublicId})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		Song: song,
	}, nil
}

type PlayPlaylistParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string `json:"playlist_public_id"`
}

func (a *Actions) PlayPlaylist(params PlayPlaylistParams) (PlaySongPayload, error) {
	err := a.app.PlayPlaylist(params.Account.Id, params.ClientHash, params.PlaylistPublicId)
	if err != nil {
		return PlaySongPayload{}, err
	}

	song, err := a.app.GetCurrentPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:        params.Account.Id,
		ClientHash:       params.ClientHash,
		SongPublicId:     song.PublicId,
		PlaylistPublicId: params.PlaylistPublicId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		Song: mapModelToActionsSong(song),
	}, nil
}

type PlaySongFromPlaylistParams struct {
	ActionContext    `json:"-"`
	SongPublicId     string `json:"song_public_id"`
	PlaylistPublicId string `json:"playlist_public_id"`
}

func (a *Actions) PlaySongFromPlaylist(params PlaySongFromPlaylistParams) (PlaySongPayload, error) {
	err := a.app.PlaySongFromPlaylist(params.Account.Id, params.ClientHash, params.SongPublicId, params.PlaylistPublicId)
	if err != nil {
		return PlaySongPayload{}, err
	}

	song, err := a.app.GetCurrentPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:        params.Account.Id,
		ClientHash:       params.ClientHash,
		SongPublicId:     params.SongPublicId,
		PlaylistPublicId: params.PlaylistPublicId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		Song: mapModelToActionsSong(song),
	}, nil
}

type PlaySongFromFavoritesParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) PlaySongFromFavorites(params PlaySongFromFavoritesParams) (PlaySongPayload, error) {
	err := a.app.PlaySongFromFavorites(params.Account.Id, params.ClientHash, params.SongPublicId)
	if err != nil {
		return PlaySongPayload{}, err
	}

	song, err := a.app.GetCurrentPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:    params.Account.Id,
		ClientHash:   params.ClientHash,
		SongPublicId: params.SongPublicId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		Song: mapModelToActionsSong(song),
	}, nil
}

type PlaySongFromQueueParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) PlaySongFromQueue(params PlaySongFromQueueParams) (PlaySongPayload, error) {
	err := a.app.PlaySongFromQueue(params.Account.Id, params.ClientHash, params.SongPublicId)
	if err != nil {
		return PlaySongPayload{}, err
	}

	song, err := a.app.GetCurrentPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return PlaySongPayload{}, err
	}

	err = a.eventhub.Publish(events.SongPlayed{
		AccountId:    params.Account.Id,
		ClientHash:   params.ClientHash,
		SongPublicId: params.SongPublicId,
	})
	if err != nil {
		return PlaySongPayload{}, err
	}

	return PlaySongPayload{
		Song: mapModelToActionsSong(song),
	}, nil
}

func (a *Actions) SaveSongsMetadataFromYouTube(songs []Song) error {
	for _, newSong := range songs {
		_, _ = a.app.CreateSong(models.Song{
			PublicId:        newSong.PublicId,
			Title:           newSong.Title,
			Artist:          newSong.Artist,
			ThumbnailUrl:    fmt.Sprintf("%s/pix/%s.webp", config.CdnAddress(), newSong.PublicId),
			RealDuration:    newSong.Duration,
			FullyDownloaded: false,
		})
	}

	return nil
}

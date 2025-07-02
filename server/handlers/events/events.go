package events

import (
	"dankmuzikk/actions"
	"dankmuzikk/app/models"
	"dankmuzikk/evy/events"
)

type EventHandlers struct {
	usecases *actions.Actions
}

func New(usecases *actions.Actions) *EventHandlers {
	return &EventHandlers{
		usecases: usecases,
	}
}

func (e *EventHandlers) HandleAddSongToHistory(event events.SongPlayed) error {
	return e.usecases.AddSongToHistory(event.SongPublicId, uint(event.AccountId))
}

func (e *EventHandlers) HandleDownloadSongOnPlay(event events.SongPlayed) error {
	return e.usecases.DownloadYouTubeSong(event.SongPublicId)
}

func (e *EventHandlers) HandleIncrementSongPlaysInPlaylist(event events.SongPlayed) error {
	return e.usecases.IncrementSongPlaysInPlaylist(event.SongPublicId, event.PlaylistPublicId, uint(event.AccountId))
}

func (e *EventHandlers) HandleMarkSongAsDownloaded(event events.SongDownloaded) error {
	return e.usecases.MarkSongAsDownloaded(event.SongPublicId)
}

func (e *EventHandlers) HandleDownloadSongOnAddingToPlaylist(event events.SongAddedToPlaylist) error {
	return e.usecases.DownloadYouTubeSong(event.SongPublicId)
}

func (e *EventHandlers) HandleIncrementPlaylistSongsCount(event events.SongAddedToPlaylist) error {
	return e.usecases.IncrementSongsCountForPlaylist(event.PlaylistPubId, event.AccountId)
}

func (e *EventHandlers) HandleDecrementPlaylistSongsCount(event events.SongRemovedFromPlaylist) error {
	return e.usecases.DecrementSongsCountForPlaylist(event.PlaylistPubId, event.AccountId)
}

func (e *EventHandlers) HandleSaveSongsMetadataOnSearchBatch(event events.SongsSearched) error {
	songs := make([]actions.Song, 0, len(event.Songs))
	for _, newSong := range event.Songs {
		songs = append(songs, actions.Song{
			PublicId:     newSong.YouTubeId,
			Title:        newSong.Title,
			Artist:       newSong.Artist,
			ThumbnailUrl: newSong.ThumbnailUrl,
			Duration:     newSong.Duration,
		})
	}

	return e.usecases.SaveSongsMetadataFromYouTube(songs)
}

func (e *EventHandlers) HandleDeletePlaylistArchive(event events.PlaylistDownloaded) error {
	return e.usecases.DeletePlaylistArchive(event)
}

func (e *EventHandlers) HandleDownloadSongOnFavorite(event events.SongAddedToFavorites) error {
	return e.usecases.DownloadYouTubeSong(event.SongPublicId)
}

func (e *EventHandlers) HandleAddSongToQueue(event events.SongPlayed) error {
	var err error
	ctx := actions.ActionContext{
		Account: models.Account{
			Id: uint(event.AccountId),
		},
		AccountId: event.AccountId,
	}
	switch event.EntryPoint {
	case events.SingleSongEntryPoint:
		err = e.usecases.AddSongToNewQueue(actions.AddSongToNewQueueParams{
			ActionContext: ctx,
			SongPublicId:  event.SongPublicId,
		})
	case events.PlayPlaylistEntryPoint:
		err = e.usecases.PlayPlaylist(actions.PlayPlaylistParams{
			ActionContext:    ctx,
			PlaylistPublicId: event.PlaylistPublicId,
		})
	case events.FromPlaylistEntryPoint:
		err = e.usecases.PlaySongFromPlaylist(actions.PlaySongFromPlaylistParams{
			ActionContext:    ctx,
			SongPublicId:     event.SongPublicId,
			PlaylistPublicId: event.PlaylistPublicId,
		})
	case events.FavoriteSongEntryPoint:
		// TODO: implement this lol
	}

	return err
}

func (e *EventHandlers) IncrementSongPlaysInPlaylist() {}

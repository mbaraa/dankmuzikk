package events

import (
	"dankmuzikk/actions"
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
	return e.usecases.AddSongToHistory(event.SongYtId, event.ProfileId)
}

func (e *EventHandlers) HandleDownloadSongOnPlay(event events.SongPlayed) error {
	return e.usecases.DownloadYouTubeSong(event.SongYtId)
}

func (e *EventHandlers) HandleIncrementSongPlaysInPlaylist(event events.SongPlayed) error {
	return e.usecases.IncrementSongPlaysInPlaylist(event.SongYtId, event.PlaylistPubId, event.ProfileId)
}

func (e *EventHandlers) HandleMarkSongAsDownloaded(event events.SongDownloaded) error {
	return e.usecases.MarkSongAsDownloaded(event.SongYtId)
}

func (e *EventHandlers) HandleDownloadSongOnAddingToPlaylist(event events.SongAddedToPlaylist) error {
	return e.usecases.DownloadYouTubeSong(event.SongYtId)
}

func (e *EventHandlers) HandleIncrementPlaylistSongsCount(event events.SongAddedToPlaylist) error {
	return nil
}

func (e *EventHandlers) HandleDecrementPlaylistSongsCount(event events.SongRemovedFromPlaylist) error {
	return nil
}

func (e *EventHandlers) IncrementSongPlaysInPlaylist() {}

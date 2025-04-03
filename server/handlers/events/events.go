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
	return e.usecases.AddSongToHistory(event.SongPublicId, event.AccountId)
}

func (e *EventHandlers) HandleDownloadSongOnPlay(event events.SongPlayed) error {
	return e.usecases.DownloadYouTubeSong(event.SongPublicId)
}

func (e *EventHandlers) HandleIncrementSongPlaysInPlaylist(event events.SongPlayed) error {
	return e.usecases.IncrementSongPlaysInPlaylist(event.SongPublicId, event.PlaylistPublicId, event.AccountId)
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

func (e *EventHandlers) IncrementSongPlaysInPlaylist() {}

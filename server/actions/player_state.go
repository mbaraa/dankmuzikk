package actions

import (
	"dankmuzikk/app/models"
	"dankmuzikk/evy/events"
	"regexp"
)

type PlayerState struct {
	Shuffled         bool   `json:"shuffled"`
	CurrentSongIndex int    `json:"current_song_index"`
	LoopMode         string `json:"loop_mode"`
	Songs            []Song `json:"songs"`
}

type GetPlayerStateParams struct {
	ActionContext `json:"-"`
}

type GetPlayerStatePayload struct {
	PlayerState PlayerState `json:"player_state"`
}

func (a *Actions) GetPlayerState(params GetPlayerStateParams) (GetPlayerStatePayload, error) {
	state, err := a.app.GetPlayerState(params.Account.Id, params.ClientHash)
	if err != nil {
		return GetPlayerStatePayload{}, err
	}

	songs := make([]Song, 0, len(state.Songs))
	for _, song := range state.Songs {
		songs = append(songs, mapModelToActionsSong(song))
	}

	return GetPlayerStatePayload{
		PlayerState: PlayerState{
			Shuffled:         state.Shuffled,
			CurrentSongIndex: state.CurrentSongIndex,
			LoopMode:         string(state.LoopMode),
			Songs:            songs,
		},
	}, nil
}

type AddSongToNewQueueParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) AddSongToNewQueue(params AddSongToNewQueueParams) error {
	return a.app.CreateSongsQueue(params.Account.Id, params.ClientHash, []string{params.SongPublicId})
}

type AddSongToQueueAtLastParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) AddSongToQueueAtLast(params AddSongToQueueAtLastParams) error {
	return a.app.AddSongToQueue(params.Account.Id, params.ClientHash, params.SongPublicId)
}

type AddSongToQueueNextParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) AddSongToQueueNext(params AddSongToQueueNextParams) error {
	return a.app.AddSongToQueueAfterCurrentSong(params.Account.Id, params.ClientHash, params.SongPublicId)
}

type AddPlaylistToQueueAtLastParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string `json:"song_public_id"`
}

func (a *Actions) AddPlaylistToQueueAtLast(params AddPlaylistToQueueAtLastParams) error {
	return a.app.AddPlaylistToQueue(params.Account.Id, params.ClientHash, params.PlaylistPublicId)
}

type AddPlaylistToQueueNextParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string `json:"song_public_id"`
}

func (a *Actions) AddPlaylistToQueueNext(params AddPlaylistToQueueNextParams) error {
	return a.app.AddPlaylistToQueueAfterCurrentSong(params.Account.Id, params.ClientHash, params.PlaylistPublicId)
}

type RemoveSongFromQueueParams struct {
	ActionContext `json:"-"`
	SongIndex     int `json:"song_index"`
}

func (a *Actions) RemoveSongFromQueue(params RemoveSongFromQueueParams) error {
	return a.app.RemoveSongFromQueue(params.Account.Id, params.ClientHash, params.SongIndex)
}

type PlayPlaylistParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string `json:"playlist_public_id"`
}

func (a *Actions) PlayPlaylist(params PlayPlaylistParams) error {
	return a.app.PlayPlaylist(params.Account.Id, params.ClientHash, params.PlaylistPublicId)
}

type PlaySongFromPlaylistParams struct {
	ActionContext    `json:"-"`
	SongPublicId     string `json:"song_public_id"`
	PlaylistPublicId string `json:"playlist_public_id"`
}

func (a *Actions) PlaySongFromPlaylist(params PlaySongFromPlaylistParams) error {
	return a.app.PlaySongFromPlaylist(params.Account.Id, params.ClientHash, params.SongPublicId, params.PlaylistPublicId)
}

type PlaySongFromFavoritesParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) PlaySongFromFavorites(params PlaySongFromFavoritesParams) error {
	return a.app.PlaySongFromFavorites(params.Account.Id, params.ClientHash, params.SongPublicId)
}

type PlaySongFromQueueParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) PlaySongFromQueue(params PlaySongFromQueueParams) error {
	return a.app.PlaySongFromQueue(params.Account.Id, params.ClientHash, params.SongPublicId)
}

type SetShuffleOnParams struct {
	ActionContext `json:"-"`
}

func (a *Actions) SetShuffleOn(params SetShuffleOnParams) error {
	return a.app.SetShuffledOn(params.Account.Id, params.ClientHash)
}

type SetShuffleOffParams struct {
	ActionContext `json:"-"`
}

func (a *Actions) SetShuffleOff(params SetShuffleOffParams) error {
	return a.app.SetShuffledOff(params.Account.Id, params.ClientHash)
}

type SetLoopOffParams struct {
	ActionContext `json:"-"`
}

func (a *Actions) SetLoopOff(params SetLoopOffParams) error {
	return a.app.SetLoopMode(params.Account.Id, params.ClientHash, models.LoopOffMode)
}

type SetLoopOnceParams struct {
	ActionContext `json:"-"`
}

func (a *Actions) SetLoopOnce(params SetLoopOnceParams) error {
	return a.app.SetLoopMode(params.Account.Id, params.ClientHash, models.LoopOnceMode)
}

type SetLoopAllParams struct {
	ActionContext `json:"-"`
}

func (a *Actions) SetLoopAll(params SetLoopAllParams) error {
	return a.app.SetLoopMode(params.Account.Id, params.ClientHash, models.LoopAllMode)
}

type GetNextSongInQueueParams struct {
	ActionContext `json:"-"`
}

type GetNextSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetNextSongInQueue(params GetNextSongInQueueParams) (GetNextSongInQueuePayload, error) {
	result, err := a.app.GetNextPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return GetNextSongInQueuePayload{}, err
	}

	event := events.SongPlayed{
		AccountId:    params.Account.Id,
		ClientHash:   params.ClientHash,
		SongPublicId: result.Song.PublicId,
		EntryPoint:   events.QueueSongEntryPoint,
	}
	err = a.eventhub.Publish(event)
	if err != nil {
		return GetNextSongInQueuePayload{}, err
	}

	// TODO: move this back to the event handler
	err = a.HandleAddSongToQueue(event)
	if err != nil {
		return GetNextSongInQueuePayload{}, err
	}

	return GetNextSongInQueuePayload{
		Song:             mapModelToActionsSong(result.Song),
		CurrentSongIndex: result.CurrentPlayingSongIndex,
		EndOfQueue:       result.EndOfQueue,
	}, nil
}

type GetPreviousSongInQueueParams struct {
	ActionContext `json:"-"`
}

type GetPreviousSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetPreviousSongInQueue(params GetPreviousSongInQueueParams) (GetPreviousSongInQueuePayload, error) {
	result, err := a.app.GetPreviousPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return GetPreviousSongInQueuePayload{}, err
	}

	event := events.SongPlayed{
		AccountId:    params.Account.Id,
		ClientHash:   params.ClientHash,
		SongPublicId: result.Song.PublicId,
		EntryPoint:   events.QueueSongEntryPoint,
	}
	err = a.eventhub.Publish(event)
	if err != nil {
		return GetPreviousSongInQueuePayload{}, err
	}

	// TODO: move this back to the event handler
	err = a.HandleAddSongToQueue(event)
	if err != nil {
		return GetPreviousSongInQueuePayload{}, err
	}

	return GetPreviousSongInQueuePayload{
		Song:             mapModelToActionsSong(result.Song),
		CurrentSongIndex: result.CurrentPlayingSongIndex,
		EndOfQueue:       result.EndOfQueue,
	}, nil
}

type GetLyricsForPlayingSongParams struct {
	ActionContext `json:"-"`
}

type GetLyricsForPlayingSongPayload struct {
	SongTitle string            `json:"song_title"`
	Lyrics    []string          `json:"lyrics"`
	Synced    map[string]string `json:"synced"`
}

var songTitleWeirdStuff = regexp.MustCompile(`(\(.*\)|\[.*\]|\{.*\}|\<.*\>)`)

func (a *Actions) GetLyricsForPlayingSong(params GetLyricsForPlayingSongParams) (GetLyricsForPlayingSongPayload, error) {
	currentSong, err := a.app.GetCurrentPlayingSong(params.Account.Id, params.ClientHash)
	if err != nil {
		return GetLyricsForPlayingSongPayload{}, err
	}

	lyrics, synced, err := a.lyrics.GetForSong(
		songTitleWeirdStuff.ReplaceAllString(currentSong.Title, ""))
	if err != nil {
		return GetLyricsForPlayingSongPayload{}, err
	}

	return GetLyricsForPlayingSongPayload{
		SongTitle: currentSong.Title,
		Lyrics:    lyrics,
		Synced:    synced,
	}, nil
}

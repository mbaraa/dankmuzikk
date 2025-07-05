package actions

type PlayerState struct {
	Shuffled         bool   `json:"shuffled"`
	CurrentSongIndex int    `json:"current_song_index"`
	LoopMode         string `json:"loop_mode"`
	Songs            []Song `json:"songs"`
}

type GetPlayerStatePayload struct {
	PlayerState PlayerState `json:"player_state"`
}

func (a *Actions) GetPlayerState(ctx ActionContext) (GetPlayerStatePayload, error) {
	return a.requests.GetPlayerState(ctx.SessionToken, ctx.ClientHash)
}

func (a *Actions) SetPlayerShuffleOn(ctx ActionContext) error {
	return a.requests.SetPlayerShuffleOn(ctx.SessionToken, ctx.ClientHash)
}

func (a *Actions) SetPlayerShuffleOff(ctx ActionContext) error {
	return a.requests.SetPlayerShuffleOff(ctx.SessionToken, ctx.ClientHash)
}

func (a *Actions) SetPlayerLoopOff(ctx ActionContext) error {
	return a.requests.SetPlayerLoopOff(ctx.SessionToken, ctx.ClientHash)
}

func (a *Actions) SetPlayerLoopOnce(ctx ActionContext) error {
	return a.requests.SetPlayerLoopOnce(ctx.SessionToken, ctx.ClientHash)
}

func (a *Actions) SetPlayerLoopAll(ctx ActionContext) error {
	return a.requests.SetPlayerLoopAll(ctx.SessionToken, ctx.ClientHash)
}

type AddSongToQueueNextParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) AddSongToQueueNext(params AddSongToQueueNextParams) error {
	return a.requests.AddSongToQueueNext(params.SessionToken, params.ClientHash, params.SongPublicId)
}

type AddSongToQueueAtLastParams struct {
	ActionContext
	SongPublicId string
}

func (a *Actions) AddSongToQueueAtLast(params AddSongToQueueNextParams) error {
	return a.requests.AddSongToQueueAtLast(params.SessionToken, params.ClientHash, params.SongPublicId)
}

type RemoveSongFromQueueParams struct {
	ActionContext
	SongIndex int
}

func (a *Actions) RemoveSongFromQueue(params RemoveSongFromQueueParams) error {
	return a.requests.RemoveSongFromQueue(params.SessionToken, params.ClientHash, params.SongIndex)
}

type AddPlaylistToQueueNextParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) AddPlaylistToQueueNext(params AddPlaylistToQueueNextParams) error {
	return a.requests.AddPlaylistToQueueNext(params.SessionToken, params.ClientHash, params.PlaylistPublicId)
}

type AddPlaylistToQueueAtLastParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) AddPlaylistToQueueAtLast(params AddPlaylistToQueueAtLastParams) error {
	return a.requests.AddPlaylistToQueueAtLast(params.SessionToken, params.ClientHash, params.PlaylistPublicId)
}

type GetNextSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetNextSongInQueue(ctx ActionContext) (GetNextSongInQueuePayload, error) {
	return a.requests.GetNextSongInQueue(ctx.SessionToken, ctx.ClientHash)
}

type GetPreviousSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetPreviousSongInQueue(ctx ActionContext) (GetPreviousSongInQueuePayload, error) {
	return a.requests.GetPreviousSongInQueue(ctx.SessionToken, ctx.ClientHash)
}

func (a *Actions) GetPlayingSongLyrics(ctx ActionContext) (GetLyricsForSongPayload, error) {
	return a.requests.GetPlayingSongLyrics(ctx.SessionToken, ctx.ClientHash)
}

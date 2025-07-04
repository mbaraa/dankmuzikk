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

func (a *Actions) GetPlayerState(sessionToken, clientHash string) (GetPlayerStatePayload, error) {
	return a.requests.GetPlayerState(sessionToken, clientHash)
}

func (a *Actions) SetPlayerShuffleOn(sessionToken, clientHash string) error {
	return a.requests.SetPlayerShuffleOn(sessionToken, clientHash)
}

func (a *Actions) SetPlayerShuffleOff(sessionToken, clientHash string) error {
	return a.requests.SetPlayerShuffleOff(sessionToken, clientHash)
}

func (a *Actions) SetPlayerLoopOff(sessionToken, clientHash string) error {
	return a.requests.SetPlayerLoopOff(sessionToken, clientHash)
}

func (a *Actions) SetPlayerLoopOnce(sessionToken, clientHash string) error {
	return a.requests.SetPlayerLoopOnce(sessionToken, clientHash)
}

func (a *Actions) SetPlayerLoopAll(sessionToken, clientHash string) error {
	return a.requests.SetPlayerLoopAll(sessionToken, clientHash)
}

func (a *Actions) AddSongToQueueNext(sessionToken, clientHash, songPublicId string) error {
	return a.requests.AddSongToQueueNext(sessionToken, clientHash, songPublicId)
}

func (a *Actions) AddSongToQueueAtLast(sessionToken, clientHash, songPublicId string) error {
	return a.requests.AddSongToQueueAtLast(sessionToken, clientHash, songPublicId)
}

func (a *Actions) AddPlaylistToQueueNext(sessionToken, clientHash, playlistPublicId string) error {
	return a.requests.AddPlaylistToQueueNext(sessionToken, clientHash, playlistPublicId)
}

func (a *Actions) AddPlaylistToQueueAtLast(sessionToken, clientHash, playlistPublicId string) error {
	return a.requests.AddPlaylistToQueueAtLast(sessionToken, clientHash, playlistPublicId)
}

type GetNextSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetNextSongInQueue(sessionToken, clientHash string) (GetNextSongInQueuePayload, error) {
	return a.requests.GetNextSongInQueue(sessionToken, clientHash)
}

type GetPreviousSongInQueuePayload struct {
	Song             Song `json:"song"`
	CurrentSongIndex int  `json:"current_song_index"`
	EndOfQueue       bool `json:"end_of_queue"`
}

func (a *Actions) GetPreviousSongInQueue(sessionToken, clientHash string) (GetPreviousSongInQueuePayload, error) {
	return a.requests.GetPreviousSongInQueue(sessionToken, clientHash)
}

func (a *Actions) GetPlayingSongLyrics(sessionToken, clientHash string) (GetLyricsForSongPayload, error) {
	return a.requests.GetPlayingSongLyrics(sessionToken, clientHash)
}

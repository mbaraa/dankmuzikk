package actions

type Requests interface {
	Auth(sessionToken string) error
	GetProfile(sessionToken string) (Profile, error)
	Logout(sessionToken string) error

	EmailLogin(params LoginUsingEmailParams) (LoginUsingEmailPayload, error)
	EmailSignup(params SignupUsingEmailParams) (SignupUsingEmailPayload, error)
	VerifyOtp(params VerifyOtpUsingEmailParams) (VerifyOtpUsingEmailPayload, error)
	GoogleLogin() (LoginUsingGooglePayload, error)
	GoogleFinishLogin(params FinishLoginUsingGoogleParams) (FinishLoginUsingGooglePayload, error)

	GetHistory(sessionToken string, pageIndex uint) ([]Song, error)

	CreatePlaylist(sessionToken string, playlist Playlist) (Playlist, error)
	GetPlaylist(sessionToken, playlistPublicId string) (Playlist, error)
	GetPlaylists(sessionToken string) ([]Playlist, error)
	GetAllPlaylistsForAddPopover(sessionToken string) ([]Playlist, map[string]bool, error)
	DeletePlaylist(sessionToken, playlistPublicId string) error
	ToggleJoinPlaylist(sessionToken, playlistPublicId string) (joined bool, err error)
	TogglePublicPlaylist(sessionToken, playlistPublicId string) (public bool, err error)
	DownloadPlaylist(sessionToken, playlistPublicId string) (string, error)

	GetFavorites(sessionToken string, pageIndex uint) (GetFavoritesPayload, error)
	AddSongToFavorites(sessionToken string, songPublicId string) error
	RemoveSongFromFavorites(sessionToken string, songPublicId string) error

	GetSongMetadata(sessionToken, clientHash, songPublicId string) (Song, error)
	PlaySong(sessionToken, clientHash, songPublicId, playlistPublicId string) (Song, error)
	ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (added bool, err error)
	UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (UpvoteSongInPlaylistPayload, error)
	DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (DownvoteSongInPlaylistPayload, error)
	GetSongLyrics(songPublicId string) (GetLyricsForSongPayload, error)

	GetPlayerState(sessionToken, clientHash string) (GetPlayerStatePayload, error)
	SetPlayerShuffleOn(sessionToken, clientHash string) error
	SetPlayerShuffleOff(sessionToken, clientHash string) error
	SetPlayerLoopOff(sessionToken, clientHash string) error
	SetPlayerLoopOnce(sessionToken, clientHash string) error
	SetPlayerLoopAll(sessionToken, clientHash string) error
	GetNextSongInQueue(sessionToken, clientHash string) (GetNextSongInQueuePayload, error)
	GetPreviousSongInQueue(sessionToken, clientHash string) (GetPreviousSongInQueuePayload, error)
	GetPlayingSongLyrics(sessionToken, clientHash string) (GetLyricsForSongPayload, error)
	AddSongToQueueNext(sessionToken, clientHash, songPublicId string) error
	AddSongToQueueAtLast(sessionToken, clientHash, songPublicId string) error
	AddPlaylistToQueueNext(sessionToken, clientHash, playlistPublicId string) error
	AddPlaylistToQueueAtLast(sessionToken, clientHash, playlistPublicId string) error

	SearchYouTube(query string) ([]Song, error)
	SearchYouTubeSuggestions(query string) ([]string, error)
}

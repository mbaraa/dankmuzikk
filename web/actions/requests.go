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

	GetSongMetadata(sessionToken, songPublicId string) (Song, error)
	PlaySong(sessionToken, songPublicId, playlistPublicId string) (string, error)
	ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (added bool, err error)
	UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (votesCount int, err error)
	DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (votesCount int, err error)
	GetSongLyrics(songPublicId string) (GetLyricsForSongPayload, error)

	SearchYouTube(query string) ([]Song, error)
	SearchYouTubeSuggestions(query string) ([]string, error)
}

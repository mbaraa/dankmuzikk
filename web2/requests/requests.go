package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
	"strconv"
)

type Requests struct {
}

func New() *Requests {
	return &Requests{}
}

func (r *Requests) Auth(sessionToken string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/auth",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}

func (r *Requests) GetProfile(sessionToken string) (actions.Profile, error) {
	return makeRequest[any, actions.Profile](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/profile",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
}

func (r *Requests) Logout(sessionToken string) error {
	_, err := makeRequest[any, actions.Profile](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/logout",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}

func (r *Requests) EmailLogin(params actions.LoginUsingEmailParams) (actions.LoginUsingEmailPayload, error) {
	return makeRequest[actions.LoginUsingEmailParams, actions.LoginUsingEmailPayload](makeRequestConfig[actions.LoginUsingEmailParams]{
		method:   http.MethodPost,
		endpoint: "/v1/login/email",
		body:     params,
	})
}

func (r *Requests) EmailSignup(params actions.SignupUsingEmailParams) (actions.SignupUsingEmailPayload, error) {
	return makeRequest[actions.SignupUsingEmailParams, actions.SignupUsingEmailPayload](makeRequestConfig[actions.SignupUsingEmailParams]{
		method:   http.MethodPost,
		endpoint: "/v1/signup/email",
		body:     params,
	})
}

func (r *Requests) VerifyOtp(params actions.VerifyOtpUsingEmailParams) (actions.VerifyOtpUsingEmailPayload, error) {
	return makeRequest[actions.VerifyOtpUsingEmailParams, actions.VerifyOtpUsingEmailPayload](makeRequestConfig[actions.VerifyOtpUsingEmailParams]{
		method:   http.MethodPost,
		endpoint: "/v1/verify-otp",
		body:     params,
	})
}

func (r *Requests) GoogleLogin() (actions.LoginUsingGooglePayload, error) {
	return makeRequest[any, actions.LoginUsingGooglePayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/login/google",
	})
}

func (r *Requests) GoogleFinishLogin(params actions.FinishLoginUsingGoogleParams) (actions.FinishLoginUsingGooglePayload, error) {
	return makeRequest[actions.FinishLoginUsingGoogleParams, actions.FinishLoginUsingGooglePayload](makeRequestConfig[actions.FinishLoginUsingGoogleParams]{
		method:   http.MethodPost,
		endpoint: "/v1/login/google/callback",
		body:     params,
	})
}

func (r *Requests) GetHistory(sessionToken string, pageIndex uint) ([]actions.Song, error) {
	return makeRequest[any, []actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/history",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"page": strconv.Itoa(int(pageIndex)),
		},
	})
}

type createPlaylistResponse struct {
	NewPlaylist actions.Playlist `json:"new_playlist"`
}

func (r *Requests) CreatePlaylist(sessionToken string, playlist actions.Playlist) (actions.Playlist, error) {
	resp, err := makeRequest[actions.Playlist, createPlaylistResponse](makeRequestConfig[actions.Playlist]{
		method:   http.MethodPost,
		endpoint: "/v1/playlist",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		body: playlist,
	})
	if err != nil {
		return actions.Playlist{}, err
	}

	return resp.NewPlaylist, nil
}

func (r *Requests) GetPlaylist(sessionToken, playlistPublicId string) (actions.Playlist, error) {
	return makeRequest[any, actions.Playlist](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/playlist",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"playlist-id": playlistPublicId,
		},
	})
}

func (r *Requests) GetPlaylists(sessionToken string) ([]actions.Playlist, error) {
	return makeRequest[any, []actions.Playlist](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/playlist/all",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
}

type getAllPlaylistsForAddPopoverResponse struct {
	Playlists        []actions.Playlist `json:"playlists"`
	SongsInPlaylists map[string]bool    `json:"songs_in_playlists"`
}

func (r *Requests) GetAllPlaylistsForAddPopover(sessionToken string) ([]actions.Playlist, map[string]bool, error) {
	resp, err := makeRequest[any, getAllPlaylistsForAddPopoverResponse](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/playlist/songs/mapped",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.Playlists, resp.SongsInPlaylists, nil
}

func (r *Requests) DeletePlaylist(sessionToken, playlistPublicId string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodDelete,
		endpoint: "/v1/playlist",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"playlist-id": playlistPublicId,
		},
	})

	return err
}

func (r *Requests) ToggleJoinPlaylist(sessionToken, playlistPublicId string) (joined bool, err error) {
	resp, err := makeRequest[any, map[string]bool](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/playlist/join",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return false, err
	}

	return resp["joined"], nil
}

func (r *Requests) TogglePublicPlaylist(sessionToken, playlistPublicId string) (public bool, err error) {
	resp, err := makeRequest[any, map[string]bool](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/playlist/public",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return false, err
	}

	return resp["public"], nil
}

func (r *Requests) DownloadPlaylist(sessionToken, playlistPublicId string) (string, error) {
	resp, err := makeRequest[any, map[string]string](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/playlist/zip",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return "", err
	}

	return resp["playlist_download_url"], nil
}

func (r *Requests) GetSongMetadata(sessionToken, songPublicId string) (actions.Song, error) {
	return makeRequest[any, actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song/single",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

func (r *Requests) PlaySong(sessionToken, songPublicId, playlistPublicId string) (string, error) {
	resp, err := makeRequest[any, map[string]string](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song/play",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"id":          songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return "", err
	}

	return resp["media_url"], nil
}

func (r *Requests) ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (added bool, err error) {
	resp, err := makeRequest[any, map[string]bool](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/playlist",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"song-id":     songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return false, err
	}

	return resp["added"], nil
}

func (r *Requests) UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (votesCount int, err error) {
	resp, err := makeRequest[any, int](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/playlist/upvote",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"song-id":     songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (r *Requests) DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (votesCount int, err error) {
	resp, err := makeRequest[any, int](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/song/playlist/downvote",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"song-id":     songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (r *Requests) GetSongLyrics(songPublicId string) (actions.GetLyricsForSongPayload, error) {
	return makeRequest[any, actions.GetLyricsForSongPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/song/lyrics",
		queryParams: map[string]string{
			"id": songPublicId,
		},
	})
}

func (r *Requests) SearchYouTube(query string) ([]actions.Song, error) {
	return makeRequest[any, []actions.Song](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/search",
		queryParams: map[string]string{
			"query": query,
		},
	})
}

func (r *Requests) SearchYouTubeSuggestions(query string) ([]string, error) {
	return makeRequest[any, []string](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/search/suggestions",
		queryParams: map[string]string{
			"query": query,
		},
	})
}

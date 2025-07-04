package requests

import (
	"dankmuzikk-web/actions"
	"net/http"
)

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

func (r *Requests) ToggleSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (added bool, err error) {
	resp, err := makeRequest[any, map[string]bool](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/playlist/song",
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

func (r *Requests) UpvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (actions.UpvoteSongInPlaylistPayload, error) {
	resp, err := makeRequest[any, actions.UpvoteSongInPlaylistPayload](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/playlist/song/upvote",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"song-id":     songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return actions.UpvoteSongInPlaylistPayload{}, err
	}

	return resp, nil
}

func (r *Requests) DownvoteSongInPlaylist(sessionToken, songPublicId, playlistPublicId string) (actions.DownvoteSongInPlaylistPayload, error) {
	resp, err := makeRequest[any, actions.DownvoteSongInPlaylistPayload](makeRequestConfig[any]{
		method:   http.MethodPut,
		endpoint: "/v1/playlist/song/downvote",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
		queryParams: map[string]string{
			"song-id":     songPublicId,
			"playlist-id": playlistPublicId,
		},
	})
	if err != nil {
		return actions.DownvoteSongInPlaylistPayload{}, err
	}

	return resp, nil
}

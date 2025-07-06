package actions

import (
	"dankmuzikk-web/requests"
	"net/http"
)

type PlaylistPermissions int8

const (
	VisitorPermission PlaylistPermissions = 1 << iota
	JoinerPermission
	OwnerPermission
	NonePermission PlaylistPermissions = 0
)

type Playlist struct {
	PublicId    string              `json:"public_id"`
	Title       string              `json:"title"`
	SongsCount  int                 `json:"songs_count"`
	Songs       []Song              `json:"songs"`
	IsPublic    bool                `json:"is_public"`
	Permissions PlaylistPermissions `json:"permissions"`
}

func (a *Actions) GetAllPlaylists(ctx ActionContext) ([]Playlist, error) {
	return requests.Do[any, []Playlist](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/playlist/all",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
		},
	})
}

type GetSinglePlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) GetSinglePlaylist(params GetSinglePlaylistParams) (Playlist, error) {
	return requests.Do[any, Playlist](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/playlist",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"playlist-id": params.PlaylistPublicId,
		},
	})
}

type CreatePlaylistParams struct {
	ActionContext
	Playlist Playlist
}

type createPlaylistResponse struct {
	NewPlaylist Playlist `json:"new_playlist"`
}

func (a *Actions) CreatePlaylist(params CreatePlaylistParams) (Playlist, error) {
	resp, err := requests.Do[Playlist, createPlaylistResponse](requests.Config[Playlist]{
		Method:   http.MethodPost,
		Endpoint: "/v1/playlist",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		Body: params.Playlist,
	})
	if err != nil {
		return Playlist{}, err
	}

	return resp.NewPlaylist, nil
}

type getAllPlaylistsForAddPopoverResponse struct {
	Playlists        []Playlist      `json:"playlists"`
	SongsInPlaylists map[string]bool `json:"songs_in_playlists"`
}

func (a *Actions) GetAllPlaylistsForAddPopover(ctx ActionContext) ([]Playlist, map[string]bool, error) {
	resp, err := requests.Do[any, getAllPlaylistsForAddPopoverResponse](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/playlist/songs/mapped",
		Headers: map[string]string{
			"Authorization": ctx.SessionToken,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.Playlists, resp.SongsInPlaylists, nil
}

type DeletePlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) DeletePlaylist(params DeletePlaylistParams) error {
	_, err := requests.Do[any, any](requests.Config[any]{
		Method:   http.MethodDelete,
		Endpoint: "/v1/playlist",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"playlist-id": params.PlaylistPublicId,
		},
	})

	return err
}

type ToggleJoinPlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) ToggleJoinPlaylist(params ToggleJoinPlaylistParams) (joined bool, err error) {
	resp, err := requests.Do[any, map[string]bool](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/playlist/join",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"playlist-id": params.PlaylistPublicId,
		},
	})
	if err != nil {
		return false, err
	}

	return resp["joined"], nil
}

type TogglePublicPlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) TogglePublicPlaylist(params TogglePublicPlaylistParams) (public bool, err error) {
	resp, err := requests.Do[any, map[string]bool](requests.Config[any]{
		Method:   http.MethodPut,
		Endpoint: "/v1/playlist/public",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"playlist-id": params.PlaylistPublicId,
		},
	})
	if err != nil {
		return false, err
	}

	return resp["public"], nil
}

type DownloadPlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) DownloadPlaylist(params DownloadPlaylistParams) (string, error) {
	resp, err := requests.Do[any, map[string]string](requests.Config[any]{
		Method:   http.MethodGet,
		Endpoint: "/v1/playlist/zip",
		Headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		QueryParams: map[string]string{
			"playlist-id": params.PlaylistPublicId,
		},
	})
	if err != nil {
		return "", err
	}

	return resp["playlist_download_url"], nil
}

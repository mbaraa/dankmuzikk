package playlists

import (
	"dankmuzikk-web/entities"
	"dankmuzikk-web/services/requests"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	httpClient *http.Client
}

func New() *Service {
	httpClient := &http.Client{
		Timeout: time.Second,
	}
	return &Service{
		httpClient: httpClient,
	}
}

func (p *Service) CreatePlaylist(token string, playlist entities.Playlist) error {
	// TODO: append this to the response using htmx magic :)
	_, err := requests.PostRequestAuth[entities.Playlist, any]("/v1/playlist", token, playlist)
	if err != nil {
		return err
	}

	return nil
}

func (p *Service) ToggleProfileInPlaylist(token, playlistPubId string) (joined bool, err error) {
	resp, err := requests.PutRequestAuth[map[string]string, map[string]bool](fmt.Sprintf("/v1/playlist/join?playlist-id=%s", playlistPubId), token, map[string]string{})
	if err != nil {
		return false, err
	}

	return resp["joined"], nil
}

func (p *Service) DeletePlaylist(token, playlistPubId string) error {
	return requests.DeleteRequestAuth("/v1/playlist?playlist-id="+playlistPubId, token)
}

func (p *Service) Get(token, playlistPubId string) (playlist entities.Playlist, err error) {
	return requests.GetRequestAuth[entities.Playlist](fmt.Sprintf("/v1/playlist?playlist-id=%s", playlistPubId), token)
}

func (p *Service) TogglePublic(token, playlistPubId string) (madePublic bool, err error) {
	resp, err := requests.PutRequestAuth[map[string]string, map[string]bool](fmt.Sprintf("/v1/playlist/public?playlist-id=%s", playlistPubId), token, map[string]string{})
	if err != nil {
		return false, err
	}

	return resp["public"], nil
}

func (p *Service) GetAll(token string) ([]entities.Playlist, error) {
	return requests.GetRequestAuth[[]entities.Playlist]("/v1/playlist/all", token)
}

func (p *Service) GetAllMappedForAddPopover(token string) ([]entities.Playlist, map[string]bool, error) {
	resp, err := requests.GetRequestAuth[struct {
		Playlists        []entities.Playlist `json:"playlists"`
		SongsInPlaylists map[string]bool     `json:"songs_in_playlists"`
	}]("/v1/playlist/songs/mapped", token)
	if err != nil {
		return nil, nil, err
	}

	return resp.Playlists, resp.SongsInPlaylists, nil
}

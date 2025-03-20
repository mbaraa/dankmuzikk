package songs

import (
	"dankmuzikk-web/entities"
	"dankmuzikk-web/services/requests"
	"errors"
	"fmt"
	"net/url"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) GetSong(songYtId string) (entities.Song, error) {
	return requests.GetRequest[entities.Song]("/v1/song/single?id=" + url.QueryEscape(songYtId))
}

func (s *Service) PlaySong(token, songYtId, playlistId string) (string, error) {
	resp, err := requests.GetRequestAuth[map[string]string](fmt.Sprintf("/v1/song/play?id=%s&playlist-id=%s", url.QueryEscape(songYtId), url.QueryEscape(playlistId)), token)
	if err != nil {
		return "", err
	}

	mediaUrl, ok := resp["media_url"]
	if !ok {
		return "", errors.New("missing media_url")
	}

	return mediaUrl, err
}

func (s *Service) ToggleSongInPlaylist(token, songId, playlistPubId string) (added bool, err error) {
	resp, err := requests.PutRequestAuth[map[string]string, map[string]bool](fmt.Sprintf("/v1/song/playlist?song-id=%s&playlist-id=%s", songId, playlistPubId), token, map[string]string{})
	if err != nil {
		return false, err
	}

	return resp["added"], nil
}

func (s *Service) UpvoteSong(token, songId, playlistPubId string) (int, error) {
	resp, err := requests.PutRequestAuth[map[string]string, int](fmt.Sprintf("/v1/song/playlist/upvote?song-id=%s&playlist-id=%s", songId, playlistPubId), token, map[string]string{})
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (s *Service) DownvoteSong(token, songId, playlistPubId string) (int, error) {
	resp, err := requests.PutRequestAuth[map[string]string, int](fmt.Sprintf("/v1/song/playlist/downvote?song-id=%s&playlist-id=%s", songId, playlistPubId), token, map[string]string{})
	if err != nil {
		return 0, err
	}

	return resp, nil
}

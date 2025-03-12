package songs

import (
	"dankmuzikk-web/entities"
	"dankmuzikk-web/services/requests"
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

func (s *Service) PlaySong(token, songYtId string) error {
	err := requests.GetRequestAuthNoRespBody("/v1/song/play?id="+url.QueryEscape(songYtId), token)
	return err
}

func (s *Service) ToggleSongInPlaylist(token, songId, playlistPubId string) (added bool, err error) {
	resp, err := requests.PutRequestAuth[map[string]string, map[string]bool](fmt.Sprintf("/v1/song/playlist?song-id=%s&playlist-id=%s", songId, playlistPubId), token, map[string]string{})
	if err != nil {
		return false, err
	}

	return resp["added"], nil
}

func (s *Service) IncrementSongPlays(token, songId, playlistPubId string) error {
	_, err := requests.PutRequestAuth[map[string]string, any](fmt.Sprintf("/v1/song/playlist/plays?song-id=%s&playlist-id=%s", songId, playlistPubId), token, map[string]string{})
	if err != nil {
		return err
	}

	return nil
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

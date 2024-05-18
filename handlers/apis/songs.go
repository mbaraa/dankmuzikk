package apis

import (
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/services/youtube/download"
	"errors"
	"net/http"
	"net/url"
)

type songDownloadHandler struct {
	service      *download.Service
	songsService *songs.Service
}

func NewDownloadHandler(service *download.Service, songsService *songs.Service) *songDownloadHandler {
	return &songDownloadHandler{service, songsService}
}

func (s *songDownloadHandler) HandleIncrementSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	songId := r.URL.Query().Get("song-id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.songsService.IncrementSongPlays(songId, playlistId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *songDownloadHandler) HandleDownloadSong(w http.ResponseWriter, r *http.Request) {
	song, err := s.extractSongFromQuery(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Errorln(err)
		return
	}

	err = s.service.DownloadYoutubeSong(song)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *songDownloadHandler) extractSongFromQuery(query url.Values) (entities.Song, error) {
	id := query.Get("yt_id")
	if id == "" {
		return entities.Song{}, errors.New("missing song's yt_id")
	}
	thumbnailUrl := query.Get("thumbnail_url")
	if thumbnailUrl == "" {
		return entities.Song{}, errors.New("missing song's thumbnail_url")
	}
	title := query.Get("title")
	if title == "" {
		return entities.Song{}, errors.New("missing song's title")
	}
	artist := query.Get("artist")
	if artist == "" {
		return entities.Song{}, errors.New("missing song's artist name")
	}
	duration := query.Get("duration")
	if duration == "" {
		return entities.Song{}, errors.New("missing song's duration")
	}

	return entities.Song{
		YtId:         id,
		Title:        title,
		Artist:       artist,
		ThumbnailUrl: thumbnailUrl,
		Duration:     duration,
	}, nil
}

package apis

import (
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/services/youtube/download"
	"net/http"
)

type songDownloadHandler struct {
	service download.Service
}

func NewDownloadHandler(service download.Service) *songDownloadHandler {
	return &songDownloadHandler{service}
}

func (s *songDownloadHandler) HandleDownloadSong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	thumbnailUrl := r.URL.Query().Get("thumbnailUrl")
	if thumbnailUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	title := r.URL.Query().Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	artist := r.URL.Query().Get("artist")
	if artist == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	duration := r.URL.Query().Get("duration")
	if duration == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.service.DownloadYoutubeSong(entities.SongDownloadRequest{
		Id:           id,
		ThumbnailUrl: thumbnailUrl,
		Title:        title,
		Artist:       artist,
		Duration:     duration,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"encoding/json"
	"fmt"
	"net/http"
)

type songDownloadHandler struct {
	usecases *actions.Actions
}

func NewDownloadHandler(usecases *actions.Actions) *songDownloadHandler {
	return &songDownloadHandler{
		usecases: usecases,
	}
}

func (s *songDownloadHandler) HandleIncrementSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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

	err := s.usecases.IncrementSongPlaysInPlaylist(songId, playlistId, profileId)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}

func (s *songDownloadHandler) HandleUpvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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

	votes, err := s.usecases.UpvoteSongInPlaylist(songId, playlistId, profileId)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songDownloadHandler) HandleDownvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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

	votes, err := s.usecases.DownvoteSongInPlaylist(songId, playlistId, profileId)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songDownloadHandler) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if profileIdCorrect {
		err := s.usecases.AddSongToHistory(id, profileId)
		if err != nil {
			log.Errorln(err)
		}
	}

	err := s.usecases.DownloadYouTubeSong(id)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}

func (s *songDownloadHandler) HandleGetSong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.GetSongByYouTubeId(id)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/app/models"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"encoding/json"
	"fmt"
	"net/http"
)

type songsHandler struct {
	usecases *actions.Actions
}

func NewSongsHandler(usecases *actions.Actions) *songsHandler {
	return &songsHandler{
		usecases: usecases,
	}
}

func (s *songsHandler) HandleUpvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
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
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songsHandler) HandleDownvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
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
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songsHandler) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlistId := r.URL.Query().Get("playlist-id")
	profileId, _ := r.Context().Value(auth.ProfileIdKey).(uint)

	mediaUrl, err := s.usecases.PlaySong(actions.PlaySongParams{
		Profile:       models.Profile{Id: profileId},
		SongYtId:      id,
		PlaylistPubId: playlistId,
	})
	if err != nil {
		log.Error("Playing a song failed", err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"media_url": mediaUrl,
	})
}

func (s *songsHandler) HandleGetSong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.GetSongByYouTubeId(id)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

package apis

import (
	"dankmuzikk/actions"
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
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
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

	payload, err := s.usecases.UpvoteSongInPlaylist(actions.UpvoteSongInPlaylistParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	// TODO: send the payload as is
	_, _ = w.Write([]byte(fmt.Sprint(payload.VotesCount)))
}

func (s *songsHandler) HandleDownvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
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

	payload, err := s.usecases.DownvoteSongInPlaylist(actions.DownvoteSongInPlaylistParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	// TODO: send the payload as is
	_, _ = w.Write([]byte(fmt.Sprint(payload.VotesCount)))
}

func (s *songsHandler) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	// un-authed action
	ctx, _ := parseContext(r.Context())

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlistId := r.URL.Query().Get("playlist-id")

	payload, err := s.usecases.PlaySong(actions.PlaySongParams{
		ActionContext: ctx,
		SongYtId:      id,
		PlaylistPubId: playlistId,
	})
	if err != nil {
		log.Error("Playing a song failed", err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsHandler) HandleGetSong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.GetSongByYouTubeId(actions.GetSongByYouTubeIdParams{
		SongYouTubeId: id,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsHandler) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fetchArtist := r.URL.Query().Get("with-artist")

	params := actions.GetLyricsForSongParams{
		SongPublicId: id,
	}
	var payload actions.GetLyricsForSongPayload
	var err error

	if fetchArtist == "true" {
		payload, err = s.usecases.GetLyricsForSongAndArtist(params)
		if err != nil {
			log.Error(err)
			handleErrorResponse(w, err)
			return
		}
	} else {
		payload, err = s.usecases.GetLyricsForSong(params)
		if err != nil {
			log.Error(err)
			handleErrorResponse(w, err)
			return
		}
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *songsHandler) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
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

	payload, err := p.usecases.ToggleSongInPlaylist(actions.ToggleSongInPlaylistParams{
		ActionContext:    ctx,
		SongPublicId:     songId,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

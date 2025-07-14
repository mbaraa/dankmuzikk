package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/evy/events"
	"encoding/json"
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
		ActionContext:    ctx,
		SongPublicId:     songId,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
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
		ActionContext:    ctx,
		SongPublicId:     songId,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsHandler) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	// un-authed action
	ctx, _ := parseContext(r.Context())

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.PlaySong(actions.PlaySongParams{
		ActionContext: ctx,
		SongPublicId:  id,
		EntryPoint:    events.SingleSongEntryPoint,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsHandler) HandlePlaySongFromPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	entryPoint := events.FromPlaylistEntryPoint
	id := r.URL.Query().Get("id")
	if id == "" {
		entryPoint = events.PlayPlaylistEntryPoint
	}
	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.PlaySong(actions.PlaySongParams{
		ActionContext: ctx,
		SongPublicId:  id,
		PlaylistPubId: playlistId,
		EntryPoint:    entryPoint,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsHandler) HandlePlaySongFromFavorites(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.PlaySong(actions.PlaySongParams{
		ActionContext: ctx,
		SongPublicId:  id,
		EntryPoint:    events.FavoriteSongEntryPoint,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsHandler) HandlePlaySongFromQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := s.usecases.PlaySong(actions.PlaySongParams{
		ActionContext: ctx,
		SongPublicId:  id,
		EntryPoint:    events.QueueSongEntryPoint,
	})
	if err != nil {
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

	ctx, _ := parseContext(r.Context())

	payload, err := s.usecases.GetSongByPublicId(actions.GetSongByPublicIdParams{
		SongPublicId:  id,
		ActionContext: ctx,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
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
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

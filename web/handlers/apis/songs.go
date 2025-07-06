package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/song"
	"dankmuzikk-web/views/components/status"
	"dankmuzikk-web/views/components/ui"
	"encoding/json"
	"net/http"
)

type songsApi struct {
	usecases *actions.Actions
}

func NewDownloadHandler(usecases *actions.Actions) *songsApi {
	return &songsApi{
		usecases: usecases,
	}
}

func (s *songsApi) HandleUpvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
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
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song.Vote(songId, playlistId, payload.VotesCount).Render(r.Context(), w)
}

func (s *songsApi) HandleDownvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
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
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song.Vote(songId, playlistId, payload.VotesCount).Render(r.Context(), w)
}

func (s *songsApi) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's id"))
		return
	}

	ctx, _ := parseContext(r.Context())

	payload, err := s.usecases.PlaySong(actions.PlaySongParams{
		ActionContext: ctx,
		SongPublicId:  id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsApi) HandlePlaySongFromPlaylist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's id"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing playlist's id"))
		return
	}

	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := s.usecases.PlaySongFromPlaylist(actions.PlaySongFromPlaylistParams{
		ActionContext:    ctx,
		SongPublicId:     id,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsApi) HandlePlaySongFromFavorites(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's id"))
		return
	}

	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := s.usecases.PlaySongFromFavorites(actions.PlaySongFromFavoritesParams{
		ActionContext: ctx,
		SongPublicId:  id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsApi) HandlePlaySongFromQueue(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's id"))
		return
	}

	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := s.usecases.PlaySongFromQueue(actions.PlaySongFromQueueParams{
		ActionContext: ctx,
		SongPublicId:  id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (s *songsApi) HandleGetSong(w http.ResponseWriter, r *http.Request) {
	ctx, _ := parseContext(r.Context())

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}

	song, err := s.usecases.GetSongMetadata(actions.GetSongMetadataParams{
		ActionContext: ctx,
		SongPublicId:  id,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(song)
}

func (p *songsApi) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
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

	added, err := p.usecases.ToggleSongInPlaylist(actions.ToggleSongInPlaylistParams{
		ActionContext:    ctx,
		SongPublicId:     songId,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	if added {
		ui.CheckedCheckbox().Render(r.Context(), w)
	} else {
		ui.UncheckedCheckbox().Render(r.Context(), w)
	}
}

package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/lyrics"
	"dankmuzikk-web/views/components/status"
	"dankmuzikk-web/views/components/ui"
	"encoding/json"
	"fmt"
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
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
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

	votes, err := s.usecases.UpvoteSongInPlaylist(sessionToken, songId, playlistId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songsApi) HandleDownvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
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

	votes, err := s.usecases.DownvoteSongInPlaylist(sessionToken, songId, playlistId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songsApi) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")

	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)

	mediaUrl, err := s.usecases.PlaySong(sessionToken, id, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"media_url": mediaUrl,
	})
}

func (s *songsApi) HandleGetSong(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}

	song, err := s.usecases.GetSongMetadata(sessionToken, id)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(song)
}

func (s *songsApi) HandleGetSongLyrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}

	lyricsResp, err := s.usecases.GetSongLyrics(id)
	if err != nil {
		status.BugsBunnyError("Lyrics was not found!").
			Render(r.Context(), w)
		return
	}

	_ = lyrics.Lyrics(lyricsResp.SongTitle, lyricsResp.Lyrics).
		Render(r.Context(), w)
}

func (p *songsApi) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
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

	added, err := p.usecases.ToggleSongInPlaylist(sessionToken, songId, playlistId)
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

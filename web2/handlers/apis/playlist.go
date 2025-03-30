package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/playlist"
	playlistviews "dankmuzikk-web/views/components/playlist"
	"dankmuzikk-web/views/components/status"
	"dankmuzikk-web/views/components/ui"
	"dankmuzikk-web/views/icons"
	"encoding/json"
	"net/http"
)

type playlistApi struct {
	usecases *actions.Actions
}

func NewPlaylistApi(usecases *actions.Actions) *playlistApi {
	return &playlistApi{
		usecases: usecases,
	}
}

func (p *playlistApi) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	var playlist actions.Playlist
	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	newPlaylist, err := p.usecases.CreatePlaylist(sessionToken, playlist)
	if err != nil {
		log.Errorln(err)
		return
	}

	playlistviews.Playlist(newPlaylist).
		Render(r.Context(), w)
}

func (p *playlistApi) HandleTogglePublicPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	madePublic, err := p.usecases.TogglePublicPlaylist(sessionToken, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	if madePublic {
		ui.CheckedCheckbox().Render(r.Context(), w)
	} else {
		ui.UncheckedCheckbox().Render(r.Context(), w)
	}
}

func (p *playlistApi) HandleToggleJoinPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	joined, err := p.usecases.ToggleJoinPlaylist(sessionToken, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	if joined {
		_ = icons.SadFrog().Render(r.Context(), w)
		_, _ = w.Write([]byte("<span>Leave playlist</span>"))
	} else {
		_ = icons.HappyFrog().Render(r.Context(), w)
		_, _ = w.Write([]byte("<span>Join playlist</span>"))
	}
}

func (p *playlistApi) HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlist, err := p.usecases.GetSinglePlaylist(sessionToken, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}
	_ = json.NewEncoder(w).Encode(playlist)
}

func (p *playlistApi) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.usecases.DeletePlaylist(sessionToken, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/playlists")
}

func (p *playlistApi) HandleGetPlaylistsForPopover(w http.ResponseWriter, r *http.Request) {
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

	playlists, songsInPlaylists, err := p.usecases.GetAllPlaylistsForAddPopover(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	playlist.PlaylistsSelector(songId, playlists, songsInPlaylists).
		Render(r.Context(), w)
}

func (p *playlistApi) HandleDonwnloadPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlistDownloadUrl, err := p.usecases.DownloadPlaylist(sessionToken, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("🤷‍♂️"))
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"playlist_download_url": playlistDownloadUrl,
	})
}

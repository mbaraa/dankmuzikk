package apis

import (
	"context"
	"dankmuzikk-web/config"
	"dankmuzikk-web/entities"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/services/playlists"
	"dankmuzikk-web/services/playlists/songs"
	"dankmuzikk-web/views/components/playlist"
	"dankmuzikk-web/views/components/ui"
	"dankmuzikk-web/views/icons"
	"dankmuzikk-web/views/pages"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type playlistApi struct {
	service     *playlists.Service
	songService *songs.Service
}

func NewPlaylistApi(service *playlists.Service, songService *songs.Service) *playlistApi {
	return &playlistApi{service, songService}
}

func (p *playlistApi) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var playlist entities.Playlist
	err = json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		w.Write([]byte("🤷‍♂️"))
		log.Errorln(err)
		return
	}

	err = p.service.CreatePlaylist(sessionToken.Value, playlist)
	if err != nil {
		log.Errorln(err)
		return
	}

	playlists, err := p.service.GetAll(sessionToken.Value)
	if err != nil {
		log.Errorln(err)
		w.Write([]byte("something went wrong"))
		return
	}

	pages.JustPlaylists(playlists).Render(context.Background(), w)
}

func (p *playlistApi) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
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

	added, err := p.songService.ToggleSongInPlaylist(sessionToken.Value, songId, playlistId)
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

func (p *playlistApi) HandleTogglePublicPlaylist(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	madePublic, err := p.service.TogglePublic(sessionToken.Value, playlistId)
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
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	joined, err := p.service.ToggleProfileInPlaylist(sessionToken.Value, playlistId)
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
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlist, err := p.service.Get(sessionToken.Value, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}
	_ = json.NewEncoder(w).Encode(playlist)
}

func (p *playlistApi) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.service.DeletePlaylist(sessionToken.Value, playlistId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/playlists")
}

func (p *playlistApi) HandleGetPlaylistsForPopover(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	songId := r.URL.Query().Get("song-id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlists, songsInPlaylists, err := p.service.GetAllMappedForAddPopover(sessionToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	playlist.PlaylistsSelector(songId, playlists, songsInPlaylists).
		Render(r.Context(), w)
}

func (p *playlistApi) HandleDonwnloadPlaylist(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("🤷‍♂️"))
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(http.MethodGet, config.GetRequestUrl("/v1/playlist/zip"), http.NoBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", sessionToken.Value)

	resp, err := (&http.Client{Timeout: time.Second * 20}).Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = io.Copy(w, resp.Body)
	_ = resp.Body.Close()
}

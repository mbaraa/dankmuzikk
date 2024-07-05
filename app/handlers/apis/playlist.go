package apis

import (
	"context"
	"dankmuzikk/entities"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/views/components/playlist"
	"dankmuzikk/views/components/ui"
	"dankmuzikk/views/icons"
	"dankmuzikk/views/pages"
	"encoding/json"
	"io"
	"net/http"
)

type playlistApi struct {
	service     *playlists.Service
	songService *songs.Service
}

func NewPlaylistApi(service *playlists.Service, songService *songs.Service) *playlistApi {
	return &playlistApi{service, songService}
}

func (p *playlistApi) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}

	var playlist entities.Playlist
	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		w.Write([]byte("🤷‍♂️"))
		log.Errorln(err)
		return
	}

	err = p.service.CreatePlaylist(playlist, profileId)
	if err != nil {
		w.Write([]byte("🤷‍♂️"))
		log.Errorln(err)
		return
	}

	playlists, err := p.service.GetAll(profileId)
	if err != nil {
		w.Write([]byte("🤷‍♂️"))
		log.Errorln(err)
		return
	}
	pages.JustPlaylists(playlists).Render(context.Background(), w)
}

func (p *playlistApi) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
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

	added, err := p.songService.ToggleSongInPlaylist(songId, playlistId, profileId)
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
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	madePublic, err := p.service.TogglePublic(playlistId, profileId)
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
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	joined, err := p.service.ToggleProfileInPlaylist(playlistId, profileId)
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
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlist, _, err := p.service.Get(playlistId, profileId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}
	_ = json.NewEncoder(w).Encode(playlist)
}

func (p *playlistApi) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.service.DeletePlaylist(playlistId, profileId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/playlists")
}

func (p *playlistApi) HandleGetPlaylistsForPopover(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte("🤷‍♂️"))
		return
	}

	songId := r.URL.Query().Get("song-id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlists, songsInPlaylists, err := p.service.GetAllMappedForAddPopover(profileId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	playlist.PlaylistsSelector(songId, playlists, songsInPlaylists).
		Render(r.Context(), w)
}

func (p *playlistApi) HandleDonwnloadPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("🤷‍♂️"))
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlistZip, err := p.service.Download(playlistId, profileId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("🤷‍♂️"))
		return
	}
	_, _ = io.Copy(w, playlistZip)
}

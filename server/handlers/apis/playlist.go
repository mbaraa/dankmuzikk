package apis

import (
	"dankmuzikk/app/entities"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/playlists/songs"
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
	playlistsList := entities.NewList(playlists, "")
	_ = json.NewEncoder(w).Encode(playlistsList)
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

	_ = json.NewEncoder(w).Encode(map[string]any{
		"added": added,
	})
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

	_ = json.NewEncoder(w).Encode(map[string]any{
		"public": madePublic,
	})
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

	_ = json.NewEncoder(w).Encode(map[string]any{
		"joined": joined,
	})
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

	// TODO: use this
	_ = songsInPlaylists

	playlistsList := entities.NewList(playlists, "")
	_ = json.NewEncoder(w).Encode(playlistsList)
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

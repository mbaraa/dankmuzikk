package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/log"
	"encoding/json"
	"io"
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
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var params actions.CreatePlaylistParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params.ProfileId = profileId

	payload, err := p.usecases.CreatePlaylist(params)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playlistApi) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlists, err := p.usecases.GetPlaylistsForProfile(profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(playlists)
}

func (p *playlistApi) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
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

	added, err := p.usecases.ToggleSongInPlaylist(songId, playlistId, profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"added": added,
	})
}

func (p *playlistApi) HandleTogglePublicPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	madePublic, err := p.usecases.TogglePublicPlaylist(playlistId, profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"public": madePublic,
	})
}

func (p *playlistApi) HandleToggleJoinPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	joined, err := p.usecases.ToggleProfileInPlaylist(playlistId, profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"joined": joined,
	})
}

func (p *playlistApi) HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlist, err := p.usecases.GetPlaylistByPublicId(playlistId, profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(playlist)
}

func (p *playlistApi) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.usecases.DeletePlaylist(playlistId, profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playlistApi) HandleGetPlaylistsForPopover(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlists, songsInPlaylists, err := p.usecases.GetAllPlaylistsMappedWithSongs(profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"playlists":          playlists,
		"songs_in_playlists": songsInPlaylists,
	})
}

func (p *playlistApi) HandleDonwnloadPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlistZip, err := p.usecases.DownloadPlaylist(playlistId, profileId)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
	_, _ = io.Copy(w, playlistZip)
}

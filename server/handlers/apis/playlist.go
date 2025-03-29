package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
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
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var params actions.CreatePlaylistParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params.ActionContext = ctx

	payload, err := p.usecases.CreatePlaylist(params)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playlistApi) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	playlists, err := p.usecases.GetPlaylistsForProfile(actions.GetPlaylistsForProfileParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(playlists)
}

func (p *playlistApi) HandleTogglePublicPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := p.usecases.TogglePublicPlaylist(actions.TogglePublicPlaylistParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playlistApi) HandleToggleJoinPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := p.usecases.ToggleJoinPlaylist(actions.ToggleJoinPlaylistParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playlistApi) HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlist, err := p.usecases.GetPlaylistByPublicId(actions.GetPlaylistByPublicIdParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(playlist)
}

func (p *playlistApi) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.DeletePlaylist(actions.DeletePlaylistParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playlistApi) HandleGetPlaylistsForPopover(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := p.usecases.GetAllPlaylistsMappedWithSongs(actions.GetAllPlaylistsMappedWithSongsParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playlistApi) HandleDonwnloadPlaylist(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	playlistId := r.URL.Query().Get("playlist-id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := p.usecases.DownloadPlaylist(actions.DownloadPlaylistParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

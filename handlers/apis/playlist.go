package apis

import (
	"context"
	"dankmuzikk/entities"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/playlists"
	"dankmuzikk/views/pages"
	"encoding/json"
	"net/http"
)

type playlistApi struct {
	service *playlists.Service
}

func NewPlaylistApi(service *playlists.Service) *playlistApi {
	return &playlistApi{service}
}

func (p *playlistApi) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
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

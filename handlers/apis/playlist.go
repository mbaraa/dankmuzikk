package apis

import (
	"context"
	"dankmuzikk/entities"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/views/pages"
	"encoding/json"
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

func (p *playlistApi) HandleToggleSongInPlaylist(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
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
	removeSongFromPlaylist := r.URL.Query().Get("remove")
	if removeSongFromPlaylist != "true" && removeSongFromPlaylist != "false" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	switch removeSongFromPlaylist {
	case "false":
		err = p.songService.AddSongToPlaylist(songId, playlistId)
	case "true":
		err = p.songService.RemoveSongFromPlaylist(songId, playlistId)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorln(err)
		return
	}

	// TODO: idk, this is ugly, but it works lol
	switch removeSongFromPlaylist {
	case "false":
		w.Write([]byte("<div class=\"w-[20px] h-[20px] rounded-sm border border-secondary bg-secondary\"></div>"))
	case "true":
		w.Write([]byte("<div class=\"w-[20px] h-[20px] rounded-sm border border-secondary\"></div>"))
	}
}

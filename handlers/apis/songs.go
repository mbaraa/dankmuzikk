package apis

import (
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/history"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/services/youtube/download"
	"net/http"
)

type songDownloadHandler struct {
	service        *download.Service
	songsService   *songs.Service
	historyService *history.Service
}

func NewDownloadHandler(
	service *download.Service,
	songsService *songs.Service,
	historyService *history.Service,
) *songDownloadHandler {
	return &songDownloadHandler{
		service:        service,
		songsService:   songsService,
		historyService: historyService,
	}
}

func (s *songDownloadHandler) HandleIncrementSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
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

	err := s.songsService.IncrementSongPlays(songId, playlistId, profileId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *songDownloadHandler) HandlePlaySong(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}

	err := s.service.DownloadYoutubeSong(id)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if !profileIdCorrect {
		log.Errorln("pashi nahui")
		return
	}
	err = s.historyService.AddSongToHistory(id, profileId)
	if err != nil {
		log.Errorln(err)
		return
	}
}

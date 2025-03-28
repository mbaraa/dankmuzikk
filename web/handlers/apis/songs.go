package apis

import (
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/services/history"
	"dankmuzikk-web/services/playlists/songs"
	"dankmuzikk-web/services/requests"
	"dankmuzikk-web/views/components/lyrics"
	"dankmuzikk-web/views/components/status"
	"encoding/json"
	"fmt"
	"net/http"
)

type songsApi struct {
	songsService   *songs.Service
	historyService *history.Service
}

func NewDownloadHandler(
	songsService *songs.Service,
	historyService *history.Service,
) *songsApi {
	return &songsApi{
		songsService:   songsService,
		historyService: historyService,
	}
}

func (s *songsApi) HandleUpvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
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

	votes, err := s.songsService.UpvoteSong(sessionToken.Value, songId, playlistId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(votes)))
}

func (s *songsApi) HandleDownvoteSongPlaysInPlaylist(w http.ResponseWriter, r *http.Request) {
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

	votes, err := s.songsService.DownvoteSong(sessionToken.Value, songId, playlistId)
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

	token := ""
	sessionToken, _ := r.Cookie(auth.SessionTokenKey)
	if sessionToken != nil {
		token = sessionToken.Value
	}

	mediaUrl, err := s.songsService.PlaySong(token, id, playlistId)
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
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}

	song, err := s.songsService.GetSong(id)
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

	lyricsResp, err := requests.GetRequest[struct {
		SongTitle string   `json:"song_title"`
		Lyrics    []string `json:"lyrics"`
	}]("/v1/song/lyrics?id=" + id)
	if err != nil {
		status.BugsBunnyError("Lyrics was not found!").
			Render(r.Context(), w)
		return
	}

	_ = lyrics.Lyrics(lyricsResp.SongTitle, lyricsResp.Lyrics).
		Render(r.Context(), w)
}

func (s *songsApi) HandleGetSongWithArtistLyrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing song's yt id"))
		return
	}
}

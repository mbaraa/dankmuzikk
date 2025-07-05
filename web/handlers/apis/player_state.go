package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/handlers/middlewares/clienthash"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/lyrics"
	"dankmuzikk-web/views/components/player"
	"dankmuzikk-web/views/components/playlist"
	"dankmuzikk-web/views/components/song"
	"dankmuzikk-web/views/components/status"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type playerStateApi struct {
	usecases *actions.Actions
}

func NewPlayerStateApi(usecases *actions.Actions) *playerStateApi {
	return &playerStateApi{
		usecases: usecases,
	}
}

func (p *playerStateApi) HandleGetPlayerState(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	payload, err := p.usecases.GetPlayerState(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPlayerSongsQueue(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	payload, err := p.usecases.GetPlayerState(sessionToken, clientHash)
	if err != nil {
		status.BugsBunnyError("No songs were found!\nMaybe play something first...").
			Render(r.Context(), w)
		return
	}

	for idx, s := range payload.PlayerState.Songs {
		song.Song(s, []string{s.AddedAt},
			[]templ.Component{
				song.RemoveFromQueue(s, idx),
				playlist.PlaylistsPopup((idx + 1), s.PublicId),
			},
			actions.Playlist{}, "queue").
			Render(r.Context(), w)
	}
}

func (p *playerStateApi) HandleSetPlayerShuffleOn(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	err := p.usecases.SetPlayerShuffleOn(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.ShuffleButton(true).Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerShuffleOff(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	err := p.usecases.SetPlayerShuffleOff(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.ShuffleButton(false).Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerLoopOff(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	err := p.usecases.SetPlayerLoopOff(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.LoopButton("off").Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerLoopOnce(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	err := p.usecases.SetPlayerLoopOnce(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.LoopButton("once").Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerLoopAll(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	err := p.usecases.SetPlayerLoopAll(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.LoopButton("all").Render(r.Context(), w)
}

func (p *playerStateApi) HandleGetNextSongInQueue(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	payload, err := p.usecases.GetNextSongInQueue(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPreviousSongInQueue(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	payload, err := p.usecases.GetPreviousSongInQueue(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPlayingSongLyrics(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	lyricsResp, err := p.usecases.GetPlayingSongLyrics(sessionToken, clientHash)
	if err != nil || len(lyricsResp.Lyrics) == 0 {
		status.BugsBunnyError("No Lyrics was found!").
			Render(r.Context(), w)
		return
	}

	_ = lyrics.Lyrics(lyricsResp.SongTitle, lyricsResp.Lyrics, lyricsResp.SyncedPairs()).
		Render(r.Context(), w)
}

func (p *playerStateApi) HandleAddSongToQueueNext(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	songId := r.URL.Query().Get("id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.usecases.AddSongToQueueNext(sessionToken, clientHash, songId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleAddSongToQueueAtLast(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	songId := r.URL.Query().Get("id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.usecases.AddSongToQueueAtLast(sessionToken, clientHash, songId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleRemoveSongFromQueue(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	songIndex, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.RemoveSongFromQueue(sessionToken, clientHash, songIndex)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleAddPlaylistToQueueNext(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	playlistId := r.URL.Query().Get("id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.usecases.AddPlaylistToQueueNext(sessionToken, clientHash, playlistId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleAddPlaylistToQueueAtLast(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	playlistId := r.URL.Query().Get("id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.usecases.AddPlaylistToQueueAtLast(sessionToken, clientHash, playlistId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

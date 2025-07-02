package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/handlers/middlewares/clienthash"
	"dankmuzikk-web/log"
	"encoding/json"
	"net/http"
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

func (p *playerStateApi) HandleSetPlayerShuffleOn(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)
	clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)

	err := p.usecases.SetPlayerShuffleOn(sessionToken, clientHash)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

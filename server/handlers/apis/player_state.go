package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
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
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	payload, err := p.usecases.GetPlayerState(actions.GetPlayerStateParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleSetShuffleOn(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	err = p.usecases.SetShuffleOn(actions.SetShuffleOnParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleSetShuffleOff(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	err = p.usecases.SetShuffleOff(actions.SetShuffleOffParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleSetLoopOff(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	err = p.usecases.SetLoopOff(actions.SetLoopOffParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleSetLoopOnce(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	err = p.usecases.SetLoopOnce(actions.SetLoopOnceParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleSetLoopAll(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	err = p.usecases.SetLoopAll(actions.SetLoopAllParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleGetNextSongInQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	payload, err := p.usecases.GetNextSongInQueue(actions.GetNextSongInQueueParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPreviousSongInQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	payload, err := p.usecases.GetPreviousSongInQueue(actions.GetPreviousSongInQueueParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleAddSongToQueueNext(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	songId := r.URL.Query().Get("id")
	if songId == "" {
		handleErrorResponse(w, ErrBadRequest{FieldName: "id"})
		return
	}

	err = p.usecases.AddSongToQueueNext(actions.AddSongToQueueNextParams{
		ActionContext: ctx,
		SongPublicId:  songId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleAddSongToQueueAtLast(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	songId := r.URL.Query().Get("id")
	if songId == "" {
		handleErrorResponse(w, ErrBadRequest{FieldName: "id"})
		return
	}

	err = p.usecases.AddSongToQueueAtLast(actions.AddSongToQueueAtLastParams{
		ActionContext: ctx,
		SongPublicId:  songId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleAddPlaylistToQueueNext(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	playlistId := r.URL.Query().Get("id")
	if playlistId == "" {
		handleErrorResponse(w, ErrBadRequest{FieldName: "id"})
		return
	}

	err = p.usecases.AddPlaylistToQueueNext(actions.AddPlaylistToQueueNextParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

func (p *playerStateApi) HandleAddPlaylistToQueueAtLast(w http.ResponseWriter, r *http.Request) {
	log.Warningln("kurwa handler")
	ctx, err := parseContext(r.Context())
	if err != nil {
		ctx, err = parseGuestContext(r.Context())
		if err != nil {
			log.Errorln(err)
			handleErrorResponse(w, err)
			return
		}
	}

	playlistId := r.URL.Query().Get("id")
	if playlistId == "" {
		handleErrorResponse(w, ErrBadRequest{FieldName: "id"})
		return
	}

	err = p.usecases.AddPlaylistToQueueAtLast(actions.AddPlaylistToQueueAtLastParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}
}

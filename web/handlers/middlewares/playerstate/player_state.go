package playerstate

import (
	"context"
	"dankmuzikk-web/actions"
	"dankmuzikk-web/handlers/middlewares/clienthash"
	"dankmuzikk-web/log"
	"net/http"
)

const PlayerStateKey = "player-state"

type mw struct {
	usecases *actions.Actions
}

func New(usecases *actions.Actions) *mw {
	return &mw{
		usecases: usecases,
	}
}

func (p *mw) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var playerState actions.PlayerState
		ctx := r.Context()
		sessionToken, err := r.Cookie("token")
		if err == nil && sessionToken != nil {
			clientHash, _ := r.Context().Value(clienthash.ClientHashKey).(string)
			playerStatePayload, err := p.usecases.GetPlayerState(actions.ActionContext{
				SessionToken: sessionToken.Value,
				ClientHash:   clientHash,
			})
			if err != nil {
				log.Errorln(err)
			}

			playerState = playerStatePayload.PlayerState
			ctx = context.WithValue(ctx, PlayerStateKey, playerState)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

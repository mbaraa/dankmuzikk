package apis

import (
	"context"
	"dankmuzikk-web/actions"
	"dankmuzikk-web/errors"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/handlers/middlewares/clienthash"
)

func parseContext(ctx context.Context) (actions.ActionContext, error) {
	sessionToken, sessionTokenCorrect := ctx.Value(auth.CtxSessionTokenKey).(string)
	if !sessionTokenCorrect {
		return actions.ActionContext{}, errors.ErrInvalidSessionToken
	}
	clientHash, clientHashCorrect := ctx.Value(clienthash.ClientHashKey).(string)
	if !clientHashCorrect {
		return actions.ActionContext{}, errors.ErrInvalidSessionToken
	}

	return actions.ActionContext{
		SessionToken: sessionToken,
		ClientHash:   clientHash,
	}, nil
}

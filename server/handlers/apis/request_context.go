package apis

import (
	"context"
	"dankmuzikk/actions"
	"dankmuzikk/app/models"
	"dankmuzikk/handlers/middlewares/auth"
)

func parseContext(ctx context.Context) (actions.ActionContext, error) {
	account, accountCorrect := ctx.Value(auth.AccountKey).(models.Account)
	if !accountCorrect {
		return actions.ActionContext{}, &ErrUnauthorized{}
	}

	return actions.ActionContext{
		Account: account,
	}, nil
}

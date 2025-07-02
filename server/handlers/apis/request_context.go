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

/* not now, not now...
func parseGuestContext(ctx context.Context) (actions.ActionContext, error) {
	guestId, guestIdCorrect := ctx.Value(auth.GuestKey).(string)
	if !guestIdCorrect {
		return actions.ActionContext{}, &ErrUnauthorized{}
	}

	return actions.ActionContext{
		AccountId: md5ToUint(guestId),
	}, nil
}

func md5ToUint(md5Str string) uint64 {
	gNum := big.NewInt(0)
	gNum.SetString(md5Str, 16)

	modResult := big.NewInt(0)
	for range 100 {
		modResult.Mod(gNum, big.NewInt(100000000069))

		if modResult.Uint64() != 0 {
			break
		}

		gNum.Add(gNum, big.NewInt(1))
	}

	return modResult.Uint64()
}
*/

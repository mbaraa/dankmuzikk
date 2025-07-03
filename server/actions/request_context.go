package actions

import "dankmuzikk/app/models"

type ActionContext struct {
	Account    models.Account
	ClientHash string
}

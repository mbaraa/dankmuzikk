package tests

import (
	"dankmuzikk/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var accounts = []models.Account{}

func initAccounts() {
	for i := 0; i < 50; i++ {
		accounts = append(accounts, models.Account{
			Email: fmt.Sprintf("%s@example.com", uuid.NewString()),
		})
	}
}

func RandomAccount() models.Account {
	return accounts[random.Intn(len(accounts))]
}

func RandomAccounts(amount int) []models.Account {
	randAccounts := make([]models.Account, amount)
	for i := 0; i < amount; i++ {
		randAccounts[i] = RandomAccount()
		random.Seed(time.Now().UnixMicro())
	}
	return randAccounts
}

func Accounts() []models.Account {
	return accounts
}

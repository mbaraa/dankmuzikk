package migrator

import (
	"dankmuzikk/db"
	"dankmuzikk/models"
)

func Migrate() error {
	dbConn, err := db.Connector()
	if err != nil {
		return err
	}

	return dbConn.Debug().AutoMigrate(
		new(models.Account),
		new(models.Profile),
		new(models.EmailVerificationCode),
		new(models.Song),
	)
}

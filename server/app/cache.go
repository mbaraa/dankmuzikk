package app

type Cache interface {
	CreateOtp(accountId uint, otp string) error
	GetOtpForAccount(id uint) (string, error)

	//CreateSongsQueue(accountId uint, initialSongs ...models.Song) error
	//GetSongsQueue(accountId uint) ([]models.Song, error)
	//AddSongToQueue(songId, accountId uint) error
	//RemoveSongFromQueue(songId, accountId uint) error
	//ResetQueue(accountId uint) error
}

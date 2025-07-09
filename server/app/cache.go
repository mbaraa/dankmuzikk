package app

import "dankmuzikk/app/models"

type Cache interface {
	CreateOtp(accountId uint, otp string) error
	GetOtpForAccount(id uint) (string, error)
}

// PlayerCache represents the glorious server state of the player for a user.
type PlayerCache interface {
	CreateSongsQueue(accountId uint, clientHash string, initialSongIds ...uint) error
	CreateSongsShuffledQueue(accountId uint, clientHash string, initialSongIds ...uint) error
	AddSongToQueue(accountId uint, clientHash string, songId uint) error
	AddSongToQueueAfterIndex(accountId uint, clientHash string, songId uint, index int) error
	AddSongToShuffledQueue(accountId uint, clientHash string, songId uint) error
	AddSongToShuffledQueueAfterIndex(accountId uint, clientHash string, songId uint, index int) error
	RemoveSongFromQueue(accountId uint, clientHash string, songIndex int) error
	RemoveSongFromShuffledQueue(accountId uint, clientHash string, songIndex int) error
	ClearQueue(accountId uint, clientHash string) error
	ClearShuffledQueue(accountId uint, clientHash string) error
	GetSongsQueue(accountId uint, clientHash string) ([]uint, error)
	GetSongsShuffledQueue(accountId uint, clientHash string) ([]uint, error)
	GetQueueLength(accountId uint, clientHash string) (uint, error)
	GetShuffledQueueLength(accountId uint, clientHash string) (uint, error)
	GetSongIdAtIndexFromQueue(accountId uint, clientHash string, index int) (uint, error)
	GetSongIdAtIndexFromShuffledQueue(accountId uint, clientHash string, index int) (uint, error)

	SetCurrentPlayingSongIndexInQueue(accountId uint, clientHash string, songIndex int) error
	SetCurrentPlayingSongIndexInShuffledQueue(accountId uint, clientHash string, songIndex int) error
	GetCurrentPlayingSongIndexInQueue(accountId uint, clientHash string) (int, error)
	GetCurrentPlayingSongIndexInShuffledQueue(accountId uint, clientHash string) (int, error)

	SetShuffled(accountId uint, clientHash string, shuffled bool) error
	GetShuffled(accountId uint, clientHash string) (bool, error)

	SetLoopMode(accountId uint, clientHash string, mode models.PlayerLoopMode) error
	GetLoopMode(accountId uint, clientHash string) (models.PlayerLoopMode, error)

	SetCurrentPlayingPlaylistInQueue(accountId uint, clientHash string, playlistId uint) error
	GetCurrentPlayingPlaylistInQueue(accountId uint, clientHash string) (uint, error)
}

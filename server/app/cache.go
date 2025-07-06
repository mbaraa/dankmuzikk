package app

import "dankmuzikk/app/models"

type Cache interface {
	CreateOtp(accountId uint, otp string) error
	GetOtpForAccount(id uint) (string, error)
}

// PlayerCache represents the glorious server state of the player for a user.
type PlayerCache interface {
	CreateSongsQueue(accountId uint, initialSongIds ...uint) error
	CreateSongsShuffledQueue(accountId uint, initialSongIds ...uint) error
	AddSongToQueue(accountId uint, songId uint) error
	AddSongToQueueAfterIndex(accountId uint, songId uint, index int) error
	AddSongToShuffledQueue(accountId uint, songId uint) error
	AddSongToShuffledQueueAfterIndex(accountId uint, songId uint, index int) error
	RemoveSongFromQueue(songIndex int, accountId uint) error
	RemoveSongFromShuffledQueue(songIndex int, accountId uint) error
	ClearQueue(accountId uint) error
	ClearShuffledQueue(accountId uint) error
	GetSongsQueue(accountId uint) ([]uint, error)
	GetSongsShuffledQueue(accountId uint) ([]uint, error)
	GetQueueLength(accountId uint) (uint, error)
	GetShuffledQueueLength(accountId uint) (uint, error)
	GetSongIdAtIndexFromQueue(accountId uint, index int) (uint, error)
	GetSongIdAtIndexFromShuffledQueue(accountId uint, index int) (uint, error)

	SetCurrentPlayingSongIndexInQueue(accountId uint, clientHash string, songIndex int) error
	SetCurrentPlayingSongIndexInShuffledQueue(accountId uint, clientHash string, songIndex int) error
	GetCurrentPlayingSongIndexInQueue(accountId uint, clientHash string) (int, error)
	GetCurrentPlayingSongIndexInShuffledQueue(accountId uint, clientHash string) (int, error)

	SetShuffled(accountId uint, shuffled bool) error
	GetShuffled(accountId uint) (bool, error)

	SetLoopMode(accountId uint, mode models.PlayerLoopMode) error
	GetLoopMode(accountId uint) (models.PlayerLoopMode, error)

	SetCurrentPlayingPlaylistInQueue(accountId uint, playlistId uint) error
	GetCurrentPlayingPlaylistInQueue(accountId uint) (uint, error)
}

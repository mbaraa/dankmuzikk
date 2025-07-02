package app

import "dankmuzikk/app/models"

type Cache interface {
	CreateOtp(accountId uint, otp string) error
	GetOtpForAccount(id uint) (string, error)
}

// PlayerCache represents the glorious server player for a user,
// or a guest, hence having account ids in [uint64] as the guests will have ids beyond [uint]'s range.
type PlayerCache interface {
	CreateSongsQueue(accountId uint64, initialSongIds ...uint) error
	CreateSongsShuffledQueue(accountId uint64, initialSongIds ...uint) error
	AddSongToQueue(accountId uint64, songId uint) error
	AddSongToQueueAfterIndex(accountId uint64, songId uint, index int) error
	AddSongToShuffledQueue(accountId uint64, songId uint) error
	AddSongToShuffledQueueAfterIndex(accountId uint64, songId uint, index int) error
	RemoveSongFromQueue(songIndex int, accountId uint64) error
	RemoveSongFromShuffledQueue(songIndex int, accountId uint64) error
	ClearQueue(accountId uint64) error
	ClearShuffledQueue(accountId uint64) error
	GetSongsQueue(accountId uint64) ([]uint, error)
	GetSongsShuffledQueue(accountId uint64) ([]uint, error)
	GetQueueLength(accountId uint64) (uint, error)
	GetShuffledQueueLength(accountId uint64) (uint, error)
	GetSongIdAtIndexFromQueue(accountId uint64, index int) (uint, error)
	GetSongIdAtIndexFromShuffledQueue(accountId uint64, index int) (uint, error)

	SetCurrentPlayingSongInedxInQueue(accountId uint64, songIndex int) error
	SetCurrentPlayingSongInedxInShuffledQueue(accountId uint64, songIndex int) error
	GetCurrentPlayingSongIndexInQueue(accountId uint64) (int, error)
	GetCurrentPlayingSongIndexInShuffledQueue(accountId uint64) (int, error)

	SetShuffled(accountId uint64, shuffled bool) error
	GetShuffled(accountId uint64) (bool, error)

	SetLoopMode(accountId uint64, mode models.PlayerLoopMode) error
	GetLoopMode(accountId uint64) (models.PlayerLoopMode, error)

	SetCurrentPlayingPlaylistInQueue(accountId uint64, playlistId uint) error
	GetCurrentPlayingPlaylistInQueue(accountId uint64) (uint, error)
}

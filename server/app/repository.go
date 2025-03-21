package app

import "dankmuzikk/app/models"

type Repository interface {
	// --------------------------------
	// User v1 stuff
	// --------------------------------

	GetAccount(id uint) (models.Account, error)
	GetAccountByEmail(email string) (models.Account, error)

	CreateProfile(profile models.Profile) (models.Profile, error)
	GetProfileForAccount(id uint) (models.Profile, error)

	CreateOtp(otp models.EmailVerificationCode) error
	GetOtpForAccount(id uint) (models.EmailVerificationCode, error)
	DeleteOtpsForAccount(id uint) error

	// --------------------------------
	// Songs v1 stuff
	// --------------------------------

	CreateSong(song models.Song) (models.Song, error)
	GetSong(id uint) (models.Song, error)
	GetSongByYouTubeId(ytId string) (models.Song, error)
	IncrementSongPlaysInPlaylist(songId, playlistPubId string, ownerId uint) error
	UpvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error)
	DownvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error)
	AddSongToHistory(songYtId string, profileId uint) error
	ToggleSongInPlaylist(songId, playlistId, ownerId uint) (added bool, err error)
	GetHistory(profileId, page uint) (models.List[models.Song], error)

	// --------------------------------
	// Playlist v1 stuff
	// --------------------------------

	GetPlaylistByPublicId(pubId string) (models.Playlist, error)
	CreatePlaylist(pl models.Playlist) (models.Playlist, error)
	AddProfileToPlaylist(plId, profileId uint) error
	RemoveProfileFromPlaylist(plId, profileId uint) error
	GetPlaylistOwners(plId uint) ([]models.PlaylistOwner, error)
	MakePlaylistPublic(id uint) error
	MakePlaylistPrivate(id uint) error
	GetPlaylistSongs(playlistId uint) (models.List[*models.Song], error)
	GetPlaylistsForProfile(ownerId uint) (models.List[models.Playlist], error)
	GetPlaylistsWithSongsForProfile(profileId uint) (models.List[models.Playlist], error)
	DeletePlaylist(id uint) error
}

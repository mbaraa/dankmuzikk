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

	// --------------------------------
	// Songs v1 stuff
	// --------------------------------

	CreateSong(song models.Song) (models.Song, error)
	GetSong(id uint) (models.Song, error)
	GetSongsByIds(ids []uint) ([]models.Song, error)
	GetSongByPublicId(pubId string) (models.Song, error)
	GetSongsByPublicIds(pubIds []string) ([]models.Song, error)
	IncrementSongPlaysInPlaylist(songId, playlistPubId string, ownerId uint) error
	UpvoteSongInPlaylist(songId, playlistPubId, ownerId uint) (int, error)
	DownvoteSongInPlaylist(songId, playlistPubId, ownerId uint) (int, error)
	AddSongToHistory(songPublicId string, accountId uint) error
	ToggleSongInPlaylist(songId, playlistId, ownerId uint) (added bool, err error)
	GetHistory(accountId, page uint) (models.List[models.Song], error)
	MarkSongAsDownloaded(songPublicId string) error
	AddSongToFavorites(songId, accountId uint) error
	RemoveSongFromFavorites(songId, accountId uint) error
	GetFavoriteSongs(accountId, page uint) (models.List[models.Song], error)

	// --------------------------------
	// Playlist v1 stuff
	// --------------------------------

	GetPlaylistByPublicId(pubId string) (models.Playlist, error)
	CreatePlaylist(pl models.Playlist) (models.Playlist, error)
	AddAccountToPlaylist(plId, accountId uint, permissions models.PlaylistPermissions) error
	RemoveAccountFromPlaylist(plId, accountId uint) error
	GetPlaylistOwners(plId uint) ([]models.PlaylistOwner, error)
	MakePlaylistPublic(id uint) error
	MakePlaylistPrivate(id uint) error
	GetPlaylistSongs(playlistId uint) (models.List[*models.Song], error)
	GetPlaylistsForAccount(accountId uint) (models.List[models.Playlist], error)
	GetPlaylistsWithSongsForAccount(account uint) (models.List[models.Playlist], error)
	DeletePlaylist(id uint) error
	IncrementSongsCountForPlaylist(id uint) error
	DecrementSongsCountForPlaylist(id uint) error
}

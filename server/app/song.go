package app

import (
	"dankmuzikk/app/models"
)

func (a *App) CreateSong(song models.Song) (models.Song, error) {
	return a.repo.CreateSong(song)
}

func (a *App) GetSongByPublicId(ytId string) (models.Song, error) {
	return a.repo.GetSongByPublicId(ytId)
}

func (a *App) IncrementSongPlaysInPlaylist(songId, playlistPubId string, ownerId uint) error {
	return a.repo.IncrementSongPlaysInPlaylist(songId, playlistPubId, ownerId)
}

func (a *App) UpvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error) {
	return a.repo.UpvoteSongInPlaylist(songId, playlistPubId, ownerId)
}

func (a *App) DownvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error) {
	return a.repo.DownvoteSongInPlaylist(songId, playlistPubId, ownerId)
}

func (a *App) AddSongToHistory(songPublicId string, accountId uint) error {
	return a.repo.AddSongToHistory(songPublicId, accountId)
}

func (a *App) ToggleSongInPlaylist(songId, playlistPubId string, ownerId uint) (added bool, err error) {
	playlist, accountPermissions, err := a.CheckAccountPlaylistAccess(ownerId, playlistPubId)
	if err != nil {
		return false, err
	}

	if accountPermissions&models.JoinerPermission == 0 && accountPermissions&models.OwnerPermission == 0 {
		return false, &ErrNotEnoughPermissionToAddSongToPlaylist{}
	}

	song, err := a.GetSongByPublicId(songId)
	if err != nil {
		return false, err
	}

	return a.repo.ToggleSongInPlaylist(song.Id, playlist.Id, ownerId)
}

func (a *App) MarkSongAsDownloaded(songPublicId string) error {
	return a.repo.MarkSongAsDownloaded(songPublicId)
}

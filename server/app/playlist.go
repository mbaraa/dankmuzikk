package app

import (
	"dankmuzikk/app/models"
	"dankmuzikk/log"
	"dankmuzikk/nanoid"
	"fmt"
)

type CreatePlaylistArgs struct {
	Title     string
	AccountId uint
}

func (a *App) CreatePlaylist(args CreatePlaylistArgs) (models.Playlist, error) {
	playlist, err := a.repo.CreatePlaylist(models.Playlist{
		PublicId: nanoid.New(),
		Title:    args.Title,
		IsPublic: false,
	})
	if err != nil {
		return models.Playlist{}, err
	}

	err = a.repo.AddAccountToPlaylist(playlist.Id, args.AccountId, models.OwnerPermission|models.JoinerPermission)
	if err != nil {
		return models.Playlist{}, err
	}

	return playlist, nil
}

func (a *App) CheckAccountPlaylistAccess(accountId uint, playlistPubId string) (models.Playlist, models.PlaylistPermissions, error) {
	playlist, err := a.repo.GetPlaylistByPublicId(playlistPubId)
	if err != nil {
		return models.Playlist{}, 0, nil
	}

	owners, err := a.repo.GetPlaylistOwners(playlist.Id)
	if err != nil {
		return models.Playlist{}, 0, err
	}

	for _, owner := range owners {
		if owner.AccountId == accountId {
			return playlist, owner.Permissions, nil
		}
	}

	if !playlist.IsPublic {
		return models.Playlist{}, 0, &ErrUnauthorizedToSeePlaylist{}
	}

	return playlist, models.VisitorPermission, nil
}

func (a *App) TogglePublicPlaylist(playlistPubId string, ownerId uint) (madePublic bool, err error) {
	playlist, accountPermissions, err := a.CheckAccountPlaylistAccess(ownerId, playlistPubId)
	if err != nil {
		return false, err
	}

	if accountPermissions&models.OwnerPermission == 0 {
		return false, &ErrNonOwnerCantChangePlaylistVisibility{}
	}

	if playlist.IsPublic {
		err = a.repo.MakePlaylistPrivate(playlist.Id)
		return false, err
	} else {
		err = a.repo.MakePlaylistPublic(playlist.Id)
		return true, err
	}
}

func (a *App) ToggleAccountInPlaylist(playlistPubId string, accountId uint) (joined bool, err error) {
	playlist, permissions, err := a.CheckAccountPlaylistAccess(accountId, playlistPubId)
	if err != nil {
		return false, err
	}

	if permissions&models.JoinerPermission == 0 {
		return true, a.repo.AddAccountToPlaylist(playlist.Id, accountId, models.JoinerPermission)
	}

	return false, a.repo.RemoveAccountFromPlaylist(playlist.Id, accountId)
}

func (a *App) GetPlaylistByPublicId(playlistPubId string, accountId uint) (models.Playlist, models.PlaylistPermissions, error) {
	playlist, permissions, err := a.CheckAccountPlaylistAccess(accountId, playlistPubId)
	if err != nil {
		return models.Playlist{}, 0, err
	}

	songs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return models.Playlist{}, 0, err
	}
	playlist.Songs = songs.Items

	return playlist, permissions, nil
}

func (a *App) DeletePlaylist(playlistPubId string, accountId uint) error {
	playlist, permissions, err := a.CheckAccountPlaylistAccess(accountId, playlistPubId)
	if err != nil {
		return err
	}
	if permissions&models.OwnerPermission == 0 {
		return &ErrNonOwnerCantDeletePlaylists{}
	}

	log.Warning("playlist", playlist.Id)

	return a.repo.DeletePlaylist(playlist.Id)
}

func (a *App) GetPlaylistsForAccount(ownerId uint) (models.List[models.Playlist], error) {
	return a.repo.GetPlaylistsForAccount(ownerId)
}

func (a *App) GetAllPlaylistsMappedWithSongs(ownerId uint) ([]models.Playlist, map[string]bool, error) {
	playlists, err := a.repo.GetPlaylistsWithSongsForAccount(ownerId)
	if err != nil {
		return nil, nil, err
	}

	mappedPlaylists := make(map[string]bool)
	for _, playlist := range playlists.Items {
		for _, song := range playlist.Songs {
			mappedPlaylists[song.YtId+"-"+playlist.PublicId] = true
		}
	}
	for i, playlist := range playlists.Items {
		for _, song := range playlist.Songs {
			if mappedPlaylists[song.YtId+"-"+playlist.PublicId] {
				continue
			}
			mappedPlaylists[fmt.Sprintf("unmapped-%d", i)] = false
		}
	}

	return playlists.Items, mappedPlaylists, nil
}

func (a *App) IncrementSongsCountForPlaylist(playlistPublicId string, accountId uint) error {
	playlist, _, err := a.GetPlaylistByPublicId(playlistPublicId, accountId)
	if err != nil {
		return err
	}

	return a.repo.IncrementSongsCountForPlaylist(playlist.Id)
}

func (a *App) DecrementSongsCountForPlaylist(playlistPublicId string, accountId uint) error {
	playlist, _, err := a.GetPlaylistByPublicId(playlistPublicId, accountId)
	if err != nil {
		return err
	}

	return a.repo.DecrementSongsCountForPlaylist(playlist.Id)
}

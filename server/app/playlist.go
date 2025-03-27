package app

import (
	"dankmuzikk/app/models"
	"dankmuzikk/log"
	"dankmuzikk/nanoid"
	"fmt"
)

type CreatePlaylistArgs struct {
	Title     string
	ProfileId uint
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

	err = a.repo.AddProfileToPlaylist(playlist.Id, args.ProfileId)
	if err != nil {
		return models.Playlist{}, err
	}

	return playlist, nil
}

func (a *App) CheckProfilePlaylistAccess(profileId uint, playlistPubId string) (models.Playlist, models.PlaylistPermissions, error) {
	playlist, err := a.repo.GetPlaylistByPublicId(playlistPubId)
	if err != nil {
		return models.Playlist{}, 0, nil
	}

	if !playlist.IsPublic {
		return models.Playlist{}, 0, &ErrUnauthorizedToSeePlaylist{}
	}

	owners, err := a.repo.GetPlaylistOwners(playlist.Id)
	if err != nil {
		return models.Playlist{}, 0, err
	}

	for _, owner := range owners {
		if owner.ProfileId == profileId {
			return playlist, owner.Permissions, nil
		}
	}

	return playlist, models.JoinerPermission | models.VisitorPermission, nil
}

func (a *App) TogglePublicPlaylist(playlistPubId string, ownerId uint) (madePublic bool, err error) {
	playlist, profilePermissions, err := a.CheckProfilePlaylistAccess(ownerId, playlistPubId)
	if err != nil {
		return false, err
	}

	if profilePermissions&models.OwnerPermission == 0 {
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

func (a *App) ToggleProfileInPlaylist(playlistPubId string, profileId uint) (joined bool, err error) {
	playlist, _, err := a.CheckProfilePlaylistAccess(profileId, playlistPubId)
	if _, ok := err.(*ErrNotFound); ok {
		return true, a.repo.AddProfileToPlaylist(playlist.Id, profileId)
	}
	if err != nil {
		return false, err
	}

	return false, a.repo.RemoveProfileFromPlaylist(playlist.Id, profileId)
}

func (a *App) GetPlaylistByPublicId(playlistPubId string, profileId uint) (models.Playlist, models.PlaylistPermissions, error) {
	playlist, permissions, err := a.CheckProfilePlaylistAccess(profileId, playlistPubId)
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

func (a *App) DeletePlaylist(playlistPubId string, profileId uint) error {
	playlist, permissions, err := a.CheckProfilePlaylistAccess(profileId, playlistPubId)
	if err != nil {
		return err
	}
	if permissions&models.OwnerPermission == 0 {
		return &ErrNonOwnerCantDeletePlaylists{}
	}

	log.Warning("playlist", playlist.Id)

	return a.repo.DeletePlaylist(playlist.Id)
}

func (a *App) GetPlaylistsForProfile(ownerId uint) (models.List[models.Playlist], error) {
	return a.repo.GetPlaylistsForProfile(ownerId)
}

func (a *App) GetAllPlaylistsMappedWithSongs(ownerId uint) ([]models.Playlist, map[string]bool, error) {
	playlists, err := a.repo.GetPlaylistsWithSongsForProfile(ownerId)
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

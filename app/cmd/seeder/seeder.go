package seeder

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/models"
	playlistspkg "dankmuzikk/services/playlists"
	"dankmuzikk/tests"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

var (
	dbConn            *gorm.DB
	accountRepo       db.UnsafeCRUDRepo[models.Account]
	profileRepo       db.UnsafeCRUDRepo[models.Profile]
	songRepo          db.UnsafeCRUDRepo[models.Song]
	playlistRepo      db.UnsafeCRUDRepo[models.Playlist]
	playlistSongsRepo db.UnsafeCRUDRepo[models.PlaylistSong]
	playlistOwnerRepo db.UnsafeCRUDRepo[models.PlaylistOwner]

	profiles  = tests.Profiles()
	songs     = tests.Songs()
	playlists = tests.Playlists()

	random = rand.New(rand.NewSource(time.Now().UnixMicro()))
)

func SeedDb() error {
	var err error

	dbConn, err = db.Connector()
	if err != nil {
		return err
	}

	accountRepo = db.NewBaseDB[models.Account](dbConn)
	profileRepo = db.NewBaseDB[models.Profile](dbConn)
	songRepo = db.NewBaseDB[models.Song](dbConn)
	playlistRepo = db.NewBaseDB[models.Playlist](dbConn)
	playlistSongsRepo = db.NewBaseDB[models.PlaylistSong](dbConn)
	playlistOwnerRepo = db.NewBaseDB[models.PlaylistOwner](dbConn)

	playlistService := playlistspkg.New(playlistRepo, playlistOwnerRepo, nil)

	pl, err := playlistService.GetAll(400)
	if err != nil {
		return err
	}
	log.Infof("%+v\n", pl)

	err = playlistService.CreatePlaylist(entities.Playlist{
		Title: "Danki Muzikki",
	}, 400)
	if err != nil {
		return err
	}

	err = playlistService.DeletePlaylist("a1a4b25f6eac4fb08222d14cadcfc7cd", 400)
	if err != nil {
		return err
	}

	return nil

	err = seedProfiles()
	if err != nil {
		return err
	}

	err = seedPlaylists()
	if err != nil {
		return err
	}

	err = addPlaylistsToProfiles()
	if err != nil {
		return err
	}

	return nil
}

func seedProfiles() error {
	for i := range profiles {
		err := profileRepo.Add(&profiles[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func seedPlaylists() error {
	for i := range playlists {
		for _, song := range playlists[i].Songs {
			err := songRepo.Add(song)
			if err != nil {
				return err
			}
		}
		err := playlistRepo.Add(&playlists[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func addPlaylistsToProfiles() error {
	for i := 0; i < len(profiles)*3; i++ {
		// ignore errors because there will be duplicates lol
		_ = playlistOwnerRepo.Add(&models.PlaylistOwner{
			ProfileId:   profiles[i%len(profiles)].Id,
			PlaylistId:  playlists[random.Intn(len(playlists))].Id,
			Permissions: models.JoinerPermission,
		})
		random.Seed(time.Now().UnixNano())
	}
	return nil
}

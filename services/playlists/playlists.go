package playlists

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"strings"

	"github.com/google/uuid"
)

// Service represents the platlist management service,
// where it fetches, adds and deletes playlists (for now)
type Service struct {
	repo               db.UnsafeCRUDRepo[models.Playlist]
	playlistOwnersRepo db.CRUDRepo[models.PlaylistOwner]
}

// New accepts a playlist repo, a playlist pwners, and returns a new instance to the playlists service.
func New(
	repo db.UnsafeCRUDRepo[models.Playlist],
	playlistOwnersRepo db.CRUDRepo[models.PlaylistOwner],
) *Service {
	return &Service{repo, playlistOwnersRepo}
}

// CreatePlaylist creates a new playlist with with provided details for the given account's profile.
// This creates a relation between profiles and playlists with the owner permission.
func (p *Service) CreatePlaylist(playlist entities.Playlist, ownerId uint) error {
	dbPlaylist := models.Playlist{
		PublicId: strings.ReplaceAll(uuid.NewString(), "-", ""),
		Title:    playlist.Title,
		IsPublic: false,
	}
	err := p.repo.Add(&dbPlaylist)
	if err != nil {
		return err
	}

	err = p.playlistOwnersRepo.Add(&models.PlaylistOwner{
		PlaylistId:  dbPlaylist.Id,
		ProfileId:   ownerId,
		Permissions: models.OwnerPermission,
	})
	if err != nil {
		return err
	}

	return nil
}

// JoinPlaylist creates a relation between profiles and playlists with write permission.
// Where only non-owners can do this, an owner literally creates the playlist ffs.
func (p *Service) JoinPlaylist(playlistPubId string, ownerId uint) error {
	dbPlaylist, err := p.repo.GetByConds("public_id = ?", playlistPubId)
	if err != nil {
		return err
	}

	err = p.playlistOwnersRepo.Add(&models.PlaylistOwner{
		PlaylistId:  dbPlaylist[0].Id,
		ProfileId:   ownerId,
		Permissions: models.WritePermission,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeletePlaylist deletes a playlist and every relation with it, that is contained songs and shared owners.
// Where only owners can do this, other permissions can just leave the playlist.
func (p *Service) DeletePlaylist(playlistPubId string, ownerId uint) error {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: ownerId,
		}).
		Where("public_id = ?", playlistPubId).
		Select("id").
		Association("Playlist").
		Find(&dbPlaylists)
	if err != nil {
		return err
	}
	if len(dbPlaylists) == 0 {
		return ErrOwnerCantLeavePlaylist
	}

	return p.
		repo.
		Delete("id = ?", dbPlaylists[0].Id)
}

// LeavePlaylist removes the relation between the given profile and the provided playlist.
// Where only non-owners can do this, since the owner can just delete the playlist, and kick everyone out :)
func (p *Service) LeavePlaylist(playlistPubId string, ownerId uint) error {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: ownerId,
		}).
		Where("public_id = ?", playlistPubId).
		Select("id").
		Association("Playlist").
		Find(&dbPlaylists)
	if err != nil {
		return err
	}
	if len(dbPlaylists) == 0 {
		return ErrNonOwnerCantDeletePlaylists
	}

	return p.
		playlistOwnersRepo.
		Delete("playlist_id = ? AND profile_id = ?", dbPlaylists[0].Id, ownerId)
}

// Get returns a full playlist (with songs) for a given profile, and an occurring error.
func (p *Service) Get(playlistPubId string, ownerId uint) (entities.Playlist, error) {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: ownerId,
		}).
		Where("public_id = ?", playlistPubId).
		Preload("Songs").
		Association("Playlist").
		Find(&dbPlaylists)
	if err != nil {
		return entities.Playlist{}, err
	}

	songs := make([]entities.Song, len(dbPlaylists[0].Songs))
	for i, song := range dbPlaylists[0].Songs {
		songs[i] = entities.Song{
			YtId:         song.YtId,
			Title:        song.Title,
			Artist:       song.Artist,
			ThumbnailUrl: song.ThumbnailUrl,
			Duration:     song.Duration,
		}
	}

	return entities.Playlist{
		PublicId:   dbPlaylists[0].PublicId,
		Title:      dbPlaylists[0].Title,
		SongsCount: dbPlaylists[0].SongsCount,
		Songs:      songs,
	}, nil
}

// GetAll returns all playlists of a profile with only meta-data (no songs), and an occurring error.
func (p *Service) GetAll(ownerId uint) ([]entities.Playlist, error) {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: ownerId,
		}).
		Association("Playlist").
		Find(&dbPlaylists)

	if err != nil {
		return nil, err
	}

	playlists := make([]entities.Playlist, len(dbPlaylists))
	for i, dbPlaylist := range dbPlaylists {
		playlists[i] = entities.Playlist{
			PublicId:   dbPlaylist.PublicId,
			Title:      dbPlaylist.Title,
			SongsCount: dbPlaylist.SongsCount,
		}
	}

	return playlists, nil
}

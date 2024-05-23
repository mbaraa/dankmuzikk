package playlists

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service represents the platlist management service,
// where it fetches, adds and deletes playlists (for now)
type Service struct {
	repo               db.UnsafeCRUDRepo[models.Playlist]
	playlistOwnersRepo db.CRUDRepo[models.PlaylistOwner]
	playlistSongsRepo  db.UnsafeCRUDRepo[models.PlaylistSong]
}

// New accepts a playlist repo, a playlist pwners, and returns a new instance to the playlists service.
func New(
	repo db.UnsafeCRUDRepo[models.Playlist],
	playlistOwnersRepo db.CRUDRepo[models.PlaylistOwner],
	playlistSongsRepo db.UnsafeCRUDRepo[models.PlaylistSong],
) *Service {
	return &Service{repo, playlistOwnersRepo, playlistSongsRepo}
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
func (p *Service) JoinPlaylist(playlistPubId string, profileId uint) error {
	dbPlaylist, err := p.repo.GetByConds("public_id = ?", playlistPubId)
	if err != nil {
		return err
	}

	err = p.playlistOwnersRepo.Add(&models.PlaylistOwner{
		PlaylistId:  dbPlaylist[0].Id,
		ProfileId:   profileId,
		Permissions: models.WritePermission,
	})
	if err != nil {
		return err
	}

	return nil
}

// LeavePlaylist removes the relation between the given profile and the provided playlist.
// Where only non-owners can do this, since the owner can just delete the playlist, and kick everyone out :)
func (p *Service) LeavePlaylist(playlistPubId string, profileId uint) error {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: profileId,
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
		Delete("playlist_id = ? AND profile_id = ?", dbPlaylists[0].Id, profileId)
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

// Get returns a full playlist (with songs) for a given profile, and an occurring error.
func (p *Service) Get(playlistPubId string, ownerId uint) (entities.Playlist, error) {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(new(models.Playlist)).
		Select("id", "is_public", "songs_count", "title", "public_id").
		Where("public_id = ?", playlistPubId).
		Find(&dbPlaylists).
		Error
	if err != nil {
		return entities.Playlist{}, err
	}
	if len(dbPlaylists) == 0 {
		return entities.Playlist{}, ErrUnauthorizedToSeePlaylist
	}
	if !dbPlaylists[0].IsPublic {
		_, err = p.playlistOwnersRepo.GetByConds("playlist_id = ? AND profile_id = ?", dbPlaylists[0].Id, ownerId)
		if err != nil {
			return entities.Playlist{}, ErrUnauthorizedToSeePlaylist
		}
	}

	gigaQuery := `SELECT yt_id, title, artist, thumbnail_url, duration, ps.created_at, ps.play_times
		FROM
			playlist_songs ps
		JOIN songs
			ON ps.song_id = songs.id
		WHERE ps.playlist_id = ?
		ORDER BY ps.created_at;`

	rows, err := p.repo.
		GetDB().
		Raw(gigaQuery, dbPlaylists[0].Id).
		Rows()
	if err != nil {
		return entities.Playlist{}, err
	}
	defer rows.Close()

	songs := make([]entities.Song, 0)
	for rows.Next() {
		var song entities.Song
		var addedAt time.Time
		err = rows.Scan(&song.YtId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.Duration, &addedAt, &song.PlayTimes)
		if err != nil {
			continue
		}
		song.AddedAt = addedAt.Format("2, January, 2006")
		songs = append(songs, song)
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
	if len(dbPlaylists) == 0 {
		return nil, ErrUnauthorizedToSeePlaylist
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

// TODO: fix this weird ass 3 return values
func (p *Service) GetAllMappedForAddPopover(ownerId uint) ([]entities.Playlist, map[string]bool, error) {
	var dbPlaylists []models.Playlist
	err := p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: ownerId,
		}).
		Preload("Songs").
		Select("id", "public_id", "title").
		Association("Playlist").
		Find(&dbPlaylists)

	if err != nil {
		return nil, nil, err
	}
	if len(dbPlaylists) == 0 {
		return nil, nil, ErrUnauthorizedToSeePlaylist
	}

	mappedPlaylists := make(map[string]bool)
	for _, playlist := range dbPlaylists {
		for _, song := range playlist.Songs {
			mappedPlaylists[song.YtId+"-"+playlist.PublicId] = true
		}
	}
	for i, playlist := range dbPlaylists {
		for _, song := range playlist.Songs {
			if mappedPlaylists[song.YtId+"-"+dbPlaylists[0].PublicId] {
				continue
			}
			mappedPlaylists[fmt.Sprintf("unmapped-%d", i)] = false
		}
	}

	playlists := make([]entities.Playlist, len(dbPlaylists))
	for i, dbPlaylist := range dbPlaylists {
		playlists[i] = entities.Playlist{
			PublicId: dbPlaylist.PublicId,
			Title:    dbPlaylist.Title,
		}
	}

	return playlists, mappedPlaylists, nil
}

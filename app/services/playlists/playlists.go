package playlists

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"dankmuzikk/services/archive"
	"dankmuzikk/services/nanoid"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Service represents the platlist management service,
// where it fetches, adds and deletes playlists (for now)
type Service struct {
	repo               db.UnsafeCRUDRepo[models.Playlist]
	playlistOwnersRepo db.CRUDRepo[models.PlaylistOwner]
	playlistSongsRepo  db.UnsafeCRUDRepo[models.PlaylistSong]
	zipService         *archive.Service
}

// New accepts a playlist repo, a playlist pwners, and returns a new instance to the playlists service.
func New(
	repo db.UnsafeCRUDRepo[models.Playlist],
	playlistOwnersRepo db.CRUDRepo[models.PlaylistOwner],
	playlistSongsRepo db.UnsafeCRUDRepo[models.PlaylistSong],
	zipService *archive.Service,
) *Service {
	return &Service{
		repo:               repo,
		playlistOwnersRepo: playlistOwnersRepo,
		playlistSongsRepo:  playlistSongsRepo,
		zipService:         zipService,
	}
}

// CreatePlaylist creates a new playlist with with provided details for the given account's profile.
// This creates a relation between profiles and playlists with the owner permission.
func (p *Service) CreatePlaylist(playlist entities.Playlist, ownerId uint) error {
	dbPlaylist := models.Playlist{
		PublicId: nanoid.Generate(),
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
		Permissions: models.OwnerPermission | models.JoinerPermission | models.VisitorPermission,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Service) ToggleProfileInPlaylist(playlistPubId string, profileId uint) (joined bool, err error) {
	playlist, err := p.repo.GetByConds("public_id = ?", playlistPubId)
	if err != nil {
		return
	}
	_, err = p.playlistOwnersRepo.GetByConds("profile_id = ? AND playlist_id = ?", profileId, playlist[0].Id)
	if errors.Is(err, db.ErrRecordNotFound) {
		err = p.playlistOwnersRepo.Add(&models.PlaylistOwner{
			ProfileId:   profileId,
			PlaylistId:  playlist[0].Id,
			Permissions: models.JoinerPermission,
		})
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		err = p.playlistOwnersRepo.Delete("profile_id = ? AND playlist_id = ?", profileId, playlist[0].Id)
		if err != nil {
			return false, err
		}
		return false, nil
	}
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
func (p *Service) Get(playlistPubId string, ownerId uint) (playlist entities.Playlist, permission models.PlaylistPermissions, err error) {
	permission = models.OwnerPermission
	var dbPlaylists []models.Playlist
	err = p.
		repo.
		GetDB().
		Model(new(models.Playlist)).
		Select("id", "is_public", "songs_count", "title", "public_id").
		Where("public_id = ?", playlistPubId).
		Find(&dbPlaylists).
		Error
	if err != nil {
		return
	}
	if len(dbPlaylists) == 0 {
		return entities.Playlist{}, models.NonePermission, ErrUnauthorizedToSeePlaylist
	}
	po, err := p.playlistOwnersRepo.GetByConds("playlist_id = ? AND profile_id = ?", dbPlaylists[0].Id, ownerId)
	if err == nil && len(po) > 0 {
		permission = po[0].Permissions
	} else {
		permission = models.VisitorPermission
	}
	if !dbPlaylists[0].IsPublic && (permission&models.JoinerPermission) == 0 {
		return entities.Playlist{}, models.NonePermission, ErrUnauthorizedToSeePlaylist
	}

	gigaQuery := `SELECT yt_id, title, artist, thumbnail_url, duration, ps.created_at, ps.play_times, ps.votes
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
		return entities.Playlist{}, models.NonePermission, err
	}
	defer rows.Close()

	songs := make([]entities.Song, 0)
	for rows.Next() {
		var song entities.Song
		var addedAt time.Time
		err = rows.Scan(&song.YtId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.Duration, &addedAt, &song.PlayTimes, &song.Votes)
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
		IsPublic:   dbPlaylists[0].IsPublic,
		Songs:      songs,
	}, permission, nil
}

// TogglePublic returns a full playlist (with songs) for a given profile, and an occurring error.
func (p *Service) TogglePublic(playlistPubId string, ownerId uint) (madePublic bool, err error) {
	var dbPlaylists []models.Playlist
	err = p.
		repo.
		GetDB().
		Model(&models.Profile{
			Id: ownerId,
		}).
		Select("id", "is_public").
		Where("public_id = ?", playlistPubId).
		Association("Playlist").
		Find(&dbPlaylists)
	if err != nil {
		return
	}

	if dbPlaylists[0].IsPublic {
		err = p.
			repo.
			GetDB().
			Model(new(models.Playlist)).
			Where("id = ?", dbPlaylists[0].Id).
			Update("is_public", false).
			Error

		return false, err
	} else {
		err = p.
			repo.
			GetDB().
			Model(new(models.Playlist)).
			Where("id = ?", dbPlaylists[0].Id).
			Update("is_public", true).
			Error

		return true, err
	}
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
			IsPublic:   dbPlaylist.IsPublic,
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

// Download zips the provided playlist,
// then returns an io.Reader with the playlist's songs, and an occurring error.
func (p *Service) Download(playlistPubId string, ownerId uint) (io.Reader, error) {
	pl, _, err := p.Get(playlistPubId, ownerId)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, len(pl.Songs))
	for i, song := range pl.Songs {
		ogFile, err := os.Open(fmt.Sprintf("%s/%s.mp3", config.Env().YouTube.MusicDir, song.YtId))
		if err != nil {
			return nil, err
		}
		newShit, err := os.OpenFile(
			fmt.Sprintf("%s/%d-%s.mp3", config.Env().YouTube.MusicDir, i+1, song.Title),
			os.O_WRONLY|os.O_CREATE, 0644,
		)
		io.Copy(newShit, ogFile)
		fileNames[i] = newShit.Name()
		_ = newShit.Close()
		_ = ogFile.Close()
	}

	zip, err := p.zipService.CreateZip()
	if err != nil {
		return nil, err
	}

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		err = zip.AddFile(file)
		if err != nil {
			return nil, err
		}
		_ = file.Close()
		_ = os.Remove(file.Name())
	}

	defer func() {
	}()

	return zip.Deflate()
}

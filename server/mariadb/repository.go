package mariadb

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/evy"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	client *gorm.DB
}

func New() (*Repository, error) {
	conn, err := dbConnector()
	if err != nil {
		return nil, err
	}

	return &Repository{
		client: conn,
	}, nil
}

func (r *Repository) GetAccount(id uint) (models.Account, error) {
	var account models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			First(&account, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) GetAccountByEmail(email string) (models.Account, error) {
	var account models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			First(&account, "email = ?", email).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) CreateProfile(profile models.Profile) (models.Profile, error) {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Profile)).
			Create(&profile).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Profile{}, &app.ErrExists{
			ResourceName: "profile",
		}
	}
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (r *Repository) GetProfileForAccount(accountId uint) (models.Profile, error) {
	var profile models.Profile

	err := tryWrapDbError(
		r.client.
			Model(new(models.Profile)).
			First(&profile, "account_id = ?", accountId).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Profile{}, &app.ErrNotFound{
			ResourceName: "profile",
		}
	}
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (r *Repository) CreateOtp(otp models.EmailVerificationCode) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.EmailVerificationCode)).
			Create(&otp).
			Error,
	)
}

func (r *Repository) GetOtpForAccount(id uint) (models.EmailVerificationCode, error) {
	var otps []models.EmailVerificationCode

	err := tryWrapDbError(
		r.client.
			Model(new(models.EmailVerificationCode)).
			Find(&otps, "account_id = ?", id).
			Error,
	)
	if _, ok := err.(ErrRecordNotFound); ok || len(otps) == 0 {
		return models.EmailVerificationCode{}, &app.ErrNotFound{
			ResourceName: "otp",
		}
	}
	if err != nil {
		return models.EmailVerificationCode{}, err
	}

	return otps[len(otps)-1], nil
}

func (r *Repository) DeleteOtpsForAccount(id uint) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.EmailVerificationCode)).
			Delete(&models.EmailVerificationCode{
				AccountId: id,
			}, "account_id = ?", id).
			Error,
	)
}

func (r *Repository) CreateSong(song models.Song) (models.Song, error) {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			Create(&song).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Song{}, &app.ErrExists{
			ResourceName: "song",
		}
	}
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (r *Repository) GetSong(id uint) (models.Song, error) {
	var song models.Song

	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			First(&song, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Song{}, &app.ErrNotFound{
			ResourceName: "song",
		}
	}
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (r *Repository) GetSongByYouTubeId(ytId string) (models.Song, error) {
	var song models.Song

	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			First(&song, "yt_id = ?", ytId).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Song{}, &app.ErrNotFound{
			ResourceName: "song",
		}
	}
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (r *Repository) IncrementSongPlaysInPlaylist(songId, playlistPubId string, ownerId uint) error {
	gigaQuery := `SELECT pl.id, s.id
		FROM
			playlist_owners po
			JOIN playlists pl
				ON po.playlist_id = pl.id
			JOIN playlist_songs ps
				ON ps.playlist_id = pl.id
			JOIN songs s
				ON ps.song_id = ps.song_id
		WHERE
			pl.public_id = ?
				AND
			s.yt_id = ?
				AND
			po.profile_id = ?
		LIMIT 1;`

	var songDbId, playlistDbId uint
	err := tryWrapDbError(
		r.client.
			Raw(gigaQuery, playlistPubId, songId, ownerId).
			Row().
			Scan(&playlistDbId, &songDbId),
	)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE playlist_songs
		SET play_times = play_times + 1
		WHERE
			playlist_id = ? AND song_id = ?;`

	err = tryWrapDbError(
		r.client.
			Exec(updateQuery, playlistDbId, songDbId).
			Error,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error) {
	gigaQuery := `SELECT pl.id, s.id
		FROM
			playlist_owners po
			JOIN playlists pl
				ON po.playlist_id = pl.id
			JOIN playlist_songs ps
				ON ps.playlist_id = pl.id
			JOIN songs s
				ON ps.song_id = ps.song_id
		WHERE
			pl.public_id = ?
				AND
			s.yt_id = ?
				AND
			po.profile_id = ?
		LIMIT 1;`

	var songDbId, playlistDbId uint
	err := tryWrapDbError(
		r.client.
			Raw(gigaQuery, playlistPubId, songId, ownerId).
			Row().
			Scan(&playlistDbId, &songDbId),
	)
	if err != nil {
		return 0, err
	}

	var voter models.PlaylistSongVoter
	err = tryWrapDbError(
		r.client.
			Model(&voter).
			First(&voter,
				"playlist_id = ? AND song_id = ? AND profile_id = ? AND vote_up = 1",
				playlistDbId,
				songDbId,
				ownerId,
			).Error,
	)
	if err == nil {
		return 0, &app.ErrUserHasAlreadyVoted{}
	}

	updateQuery := `UPDATE playlist_songs
		SET votes = votes + 1
		WHERE
			playlist_id = ? AND song_id = ?;`

	err = tryWrapDbError(
		r.client.
			Exec(updateQuery, playlistDbId, songDbId).
			Error,
	)
	if err != nil {
		return 0, err
	}

	var ps models.PlaylistSong
	err = tryWrapDbError(
		r.client.
			Model(&ps).
			First(&ps,
				"playlist_id = ? AND song_id = ?",
				playlistDbId,
				songDbId,
			).Error,
	)
	if err == nil {
		return 0, &app.ErrUserHasAlreadyVoted{}
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.PlaylistSongVoter)).
			Create(
				&models.PlaylistSongVoter{
					PlaylistId: playlistDbId,
					SongId:     songDbId,
					ProfileId:  ownerId,
					VoteUp:     true,
				}).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return ps.Votes,
			r.client.
				Exec("UPDATE playlist_song_voters SET vote_up = 1 WHERE playlist_id = ? AND song_id = ? AND profile_id = ?", playlistDbId, songDbId, ownerId).
				Error
	} else {
		return ps.Votes, err
	}
}

func (r *Repository) DownvoteSongInPlaylist(songId, playlistPubId string, ownerId uint) (int, error) {
	gigaQuery := `SELECT pl.id, s.id
		FROM
			playlist_owners po
			JOIN playlists pl
				ON po.playlist_id = pl.id
			JOIN playlist_songs ps
				ON ps.playlist_id = pl.id
			JOIN songs s
				ON ps.song_id = ps.song_id
		WHERE
			pl.public_id = ?
				AND
			s.yt_id = ?
				AND
			po.profile_id = ?
		LIMIT 1;`

	var songDbId, playlistDbId uint
	err := tryWrapDbError(
		r.client.
			Raw(gigaQuery, playlistPubId, songId, ownerId).
			Row().
			Scan(&playlistDbId, &songDbId),
	)
	if err != nil {
		return 0, err
	}

	var voter models.PlaylistSongVoter
	err = tryWrapDbError(
		r.client.
			Model(&voter).
			First(&voter,
				"playlist_id = ? AND song_id = ? AND profile_id = ? AND vote_up = 0",
				playlistDbId,
				songDbId,
				ownerId,
			).Error,
	)
	if err == nil {
		return 0, &app.ErrUserHasAlreadyVoted{}
	}

	updateQuery := `UPDATE playlist_songs
		SET votes = votes - 1
		WHERE
			playlist_id = ? AND song_id = ?;`

	err = tryWrapDbError(
		r.client.
			Exec(updateQuery, playlistDbId, songDbId).
			Error,
	)
	if err != nil {
		return 0, err
	}

	var ps models.PlaylistSong
	err = tryWrapDbError(
		r.client.
			Model(&ps).
			First(&ps,
				"playlist_id = ? AND song_id = ?",
				playlistDbId,
				songDbId,
			).Error,
	)
	if err == nil {
		return 0, &app.ErrUserHasAlreadyVoted{}
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.PlaylistSongVoter)).
			Create(
				&models.PlaylistSongVoter{
					PlaylistId: playlistDbId,
					SongId:     songDbId,
					ProfileId:  ownerId,
					VoteUp:     true,
				}).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return ps.Votes,
			r.client.
				Exec("UPDATE playlist_song_voters SET vote_up = 0 WHERE playlist_id = ? AND song_id = ? AND profile_id = ?", playlistDbId, songDbId, ownerId).
				Error
	} else {
		return ps.Votes, err
	}
}

func (r *Repository) AddSongToHistory(songYtId string, profileId uint) error {
	var song models.Song
	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			First(&song, "yt_id = ?", songYtId).
			Error,
	)
	if err != nil {
		return err
	}

	return tryWrapDbError(
		r.client.
			Model(new(models.History)).
			Create(
				&models.History{
					ProfileId: profileId,
					SongId:    song.Id,
				}).
			Error,
	)
}

func (r *Repository) ToggleSongInPlaylist(songId, playlistId, ownerId uint) (added bool, err error) {
	err = tryWrapDbError(
		r.client.
			Model(new(models.PlaylistSong)).
			First(&models.PlaylistSong{}, "playlist_id = ? AND song_id = ?", playlistId, songId).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		err = tryWrapDbError(
			r.client.
				Model(new(models.PlaylistSong)).
				Create(&models.PlaylistSong{
					PlaylistId: playlistId,
					SongId:     songId,
					Votes:      1,
				}).
				Error,
		)
		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		return false, tryWrapDbError(
			r.client.
				Model(new(models.PlaylistSong)).
				Delete(&models.PlaylistSong{
					PlaylistId: playlistId,
					SongId:     songId,
				}, "playlist_id = ? AND song_id = ?", playlistId, songId).
				Error,
		)
	}
}

func (r *Repository) GetHistory(profileId, page uint) (models.List[models.Song], error) {
	gigaQuery := fmt.Sprintf(
		`SELECT yt_id, title, artist, thumbnail_url, duration, h.created_at
		FROM
			histories h JOIN songs
		ON
				songs.id = h.song_id
		WHERE h.profile_id = ?
		ORDER BY h.created_at DESC
		LIMIT %d,%d;`,
		(page-1)*20, page*20,
	)

	rows, err := r.client.
		Raw(gigaQuery, profileId).
		Rows()
	if err != nil {
		return models.List[models.Song]{}, err
	}

	songs := make([]models.Song, 0)
	for rows.Next() {
		var song models.Song
		var addedAt time.Time
		err = rows.Scan(&song.YtId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.Duration, &addedAt)
		if err != nil {
			continue
		}
		song.AddedAt = whenDidItHappen(addedAt)
		songs = append(songs, song)
	}
	_ = rows.Close()

	return models.NewList(songs, fmt.Sprint(page+1)), nil
}

func (r *Repository) GetPlaylistSongs(playlistId uint) (models.List[*models.Song], error) {
	gigaQuery := `SELECT yt_id, title, artist, thumbnail_url, duration, ps.created_at, ps.play_times, ps.votes
		FROM
			playlist_songs ps
		JOIN songs
			ON ps.song_id = songs.id
		WHERE ps.playlist_id = ?
		ORDER BY ps.created_at;`

	rows, err := r.client.
		Raw(gigaQuery, playlistId).
		Rows()
	err = tryWrapDbError(err)
	if err != nil {
		return models.List[*models.Song]{}, err
	}

	songs := make([]*models.Song, 0)
	for rows.Next() {
		var song models.Song
		var addedAt time.Time
		err = rows.Scan(&song.YtId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.Duration, &addedAt, &song.PlayTimes, &song.Votes)
		if err != nil {
			continue
		}
		song.AddedAt = addedAt.Format("2, January, 2006")
		songs = append(songs, &song)
	}

	_ = rows.Close()

	return models.NewList(songs, ""), nil
}

func (r *Repository) MarkSongAsDownloaded(songYtId string) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			Where("yt_id = ?", songYtId).
			Update("fully_downloaded", true).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "song",
		}
	}

	return nil
}

func (r *Repository) GetPlaylistByPublicId(pubId string) (models.Playlist, error) {
	var playlist models.Playlist
	err := tryWrapDbError(
		r.client.
			Model(new(models.Playlist)).
			First(&playlist, "public_id = ?", pubId).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Playlist{}, &app.ErrNotFound{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return models.Playlist{}, err
	}

	return playlist, nil
}

func (r *Repository) CreatePlaylist(pl models.Playlist) (models.Playlist, error) {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Playlist)).
			Create(&pl).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Playlist{}, &app.ErrExists{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return models.Playlist{}, err
	}

	return pl, nil
}

func (r *Repository) AddProfileToPlaylist(plId, profileId uint) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.PlaylistOwner)).
			Create(&models.PlaylistOwner{
				PlaylistId:  plId,
				ProfileId:   profileId,
				Permissions: models.OwnerPermission | models.JoinerPermission | models.VisitorPermission,
			}).
			Error,
	)
}

func (r *Repository) RemoveProfileFromPlaylist(plId, profileId uint) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.PlaylistOwner)).
			Delete(&models.PlaylistOwner{},
				"profile_id = ? AND playlist_id = ?", profileId, plId,
			).
			Error,
	)
}

func (r *Repository) GetPlaylistOwners(plId uint) ([]models.PlaylistOwner, error) {
	var owners []models.PlaylistOwner
	err := tryWrapDbError(
		r.client.
			Model(new(models.PlaylistOwner)).
			Find(&owners, "playlist_id = ?", plId).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return nil, &app.ErrExists{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return nil, err
	}

	return owners, nil
}

func (r *Repository) MakePlaylistPublic(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Playlist)).
			Where("id = ?", id).
			Update("is_public", true).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return &app.ErrExists{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) MakePlaylistPrivate(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Playlist)).
			Where("id = ?", id).
			Update("is_public", false).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return &app.ErrExists{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPlaylistsForProfile(ownerId uint) (models.List[models.Playlist], error) {
	var playlists []models.Playlist
	err := tryWrapDbError(
		r.client.
			Model(&models.Profile{
				Id: ownerId,
			}).
			Association("Playlist").
			Find(&playlists),
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.List[models.Playlist]{}, &app.ErrNotFound{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return models.List[models.Playlist]{}, err
	}
	if len(playlists) == 0 {
		return models.List[models.Playlist]{}, &app.ErrUnauthorizedToSeePlaylist{}
	}

	return models.NewList(playlists, ""), nil
}

func (r *Repository) GetPlaylistsWithSongsForProfile(profileId uint) (models.List[models.Playlist], error) {
	var dbPlaylists []models.Playlist
	err := tryWrapDbError(
		r.client.
			Model(&models.Profile{
				Id: profileId,
			}).
			Preload("Songs").
			Select("id", "public_id", "title").
			Association("Playlist").
			Find(&dbPlaylists),
	)

	if err != nil {
		return models.List[models.Playlist]{}, err
	}
	if len(dbPlaylists) == 0 {
		return models.List[models.Playlist]{}, &app.ErrUnauthorizedToSeePlaylist{}
	}

	return models.NewList(dbPlaylists, ""), nil
}

func (r *Repository) DeletePlaylist(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Playlist)).
			Delete(&models.Playlist{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "playlist",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------
// Evy Repository
// --------------------------------

func (r *Repository) CreateEvent(e evy.EventPayload) error {
	err := tryWrapDbError(
		r.client.
			Model(new(evy.EventPayload)).
			Create(&e).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return &app.ErrExists{
			ResourceName: "event",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetEventsBatch(size int32) ([]evy.EventPayload, error) {
	var events []evy.EventPayload
	err := tryWrapDbError(
		r.client.
			Model(&evy.EventPayload{}).
			Find(&events).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "event",
		}
	}
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, &app.ErrNotFound{
			ResourceName: "event",
		}
	}

	return events, nil
}

func (r *Repository) DeleteEvent(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(evy.EventPayload)).
			Delete(&evy.EventPayload{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "event",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------
// Utils
// --------------------------------

func whenDidItHappen(t time.Time) string {
	now := time.Now().UTC()
	switch {
	case t.Day() == now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "Today"
	case t.Day()+1 == now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "Yesterday"
	case t.Day()+5 < now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "Last week"
	case t.Day() == now.Day() && t.Month()+1 == now.Month() && t.Year() == now.Year():
		return "Last month"
	default:
		return fmt.Sprintf("%s %s %s", t.Format("January"), nth(t.Day()), t.Format("2006"))
	}
}

func nth(n int) string {
	switch {
	case n%10 == 1:
		return fmt.Sprintf("%dst", n)
	case n%10 == 2:
		return fmt.Sprintf("%dnd", n)
	case n%10 == 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}

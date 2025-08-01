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

// --------------------------------
// App Repository
// --------------------------------

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

func (r *Repository) GetSongsByIds(ids []uint) ([]models.Song, error) {
	var songs []models.Song

	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			Find(&songs, "id IN ?", ids).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "song",
		}
	}
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *Repository) GetSongByPublicId(publicId string) (models.Song, error) {
	var song models.Song

	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			First(&song, "public_id = ?", publicId).
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

func (r *Repository) GetSongsByPublicIds(publicIds []string) ([]models.Song, error) {
	var songs []models.Song

	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			Find(&songs, "public_id IN ?", publicIds).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "song",
		}
	}
	if err != nil {
		return nil, err
	}

	return songs, nil
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
			s.public_id = ?
				AND
			po.account_id = ?
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

func (r *Repository) UpvoteSongInPlaylist(songId, playlistId, accountId uint) (int, error) {
	var voter models.PlaylistSongVoter
	err := tryWrapDbError(
		r.client.
			Model(&voter).
			First(&voter,
				"playlist_id = ? AND song_id = ? AND account_id = ? AND vote_up = 1",
				playlistId,
				songId,
				accountId,
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
			Exec(updateQuery, playlistId, songId).
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
				playlistId,
				songId,
			).Error,
	)
	if err != nil {
		return 0, &app.ErrUserHasAlreadyVoted{}
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.PlaylistSongVoter)).
			Create(
				&models.PlaylistSongVoter{
					PlaylistId: playlistId,
					SongId:     songId,
					AccountId:  accountId,
					VoteUp:     true,
				}).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return ps.Votes,
			r.client.
				Exec(
					"UPDATE playlist_song_voters SET vote_up = 1 WHERE playlist_id = ? AND song_id = ? AND account_id = ?",
					playlistId, songId, accountId).
				Error
	} else {
		return ps.Votes, err
	}
}

func (r *Repository) DownvoteSongInPlaylist(songId, playlistId, accountId uint) (int, error) {
	var voter models.PlaylistSongVoter
	err := tryWrapDbError(
		r.client.
			Model(&voter).
			First(&voter,
				"playlist_id = ? AND song_id = ? AND account_id = ? AND vote_up = 0",
				playlistId,
				songId,
				accountId,
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
			Exec(updateQuery, playlistId, songId).
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
				playlistId,
				songId,
			).Error,
	)
	if err != nil {
		return 0, &app.ErrUserHasAlreadyVoted{}
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.PlaylistSongVoter)).
			Create(
				&models.PlaylistSongVoter{
					PlaylistId: playlistId,
					SongId:     songId,
					AccountId:  accountId,
					VoteUp:     true,
				}).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return ps.Votes,
			r.client.
				Exec(
					"UPDATE playlist_song_voters SET vote_up = 0 WHERE playlist_id = ? AND song_id = ? AND account_id = ?",
					playlistId, songId, accountId).
				Error
	} else {
		return ps.Votes, err
	}
}

func (r *Repository) AddSongToHistory(songPublicId string, accountId uint) error {
	var song models.Song
	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			First(&song, "public_id = ?", songPublicId).
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
					AccountId: accountId,
					SongId:    song.Id,
				}).
			Error,
	)
}

func (r *Repository) GetHistory(accountId, page uint) (models.List[models.Song], error) {
	gigaQuery := fmt.Sprintf(
		`SELECT songs.id, public_id, title, artist, thumbnail_url, real_duration, h.created_at
		FROM
			histories h JOIN songs
		ON
				songs.id = h.song_id
		WHERE h.account_id = ?
		ORDER BY h.created_at DESC
		LIMIT %d,%d;`,
		(page-1)*20, page*20,
	)

	rows, err := r.client.
		Raw(gigaQuery, accountId).
		Rows()
	if err != nil {
		return models.List[models.Song]{}, err
	}

	songs := make([]models.Song, 0, 20)
	songIds := make([]uint, 0, 20)
	for rows.Next() {
		var song models.Song
		var addedAt time.Time
		err = rows.Scan(&song.Id, &song.PublicId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.RealDuration, &addedAt)
		if err != nil {
			continue
		}
		song.AddedAt = whenDidItHappen(addedAt)
		songs = append(songs, song)
		songIds = append(songIds, song.Id)
	}
	_ = rows.Close()

	rows, err = r.client.
		Raw(`SELECT song_id FROM favorite_songs WHERE account_id = ? AND song_id IN ?`, accountId, songIds).
		Rows()
	if err != nil {
		return models.List[models.Song]{}, err
	}

	songInFavorites := map[uint]bool{}
	for rows.Next() {
		var songId uint
		err = rows.Scan(&songId)
		if err != nil {
			continue
		}
		songInFavorites[songId] = true
	}

	for i := range songs {
		songs[i].Favorite = songInFavorites[songs[i].Id]
	}

	return models.NewList(songs, fmt.Sprint(page+1)), nil
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
		err = tryWrapDbError(
			r.client.
				Model(new(models.PlaylistSong)).
				Delete(&models.PlaylistSong{
					PlaylistId: playlistId,
					SongId:     songId,
				}, "playlist_id = ? AND song_id = ?", playlistId, songId).
				Error,
		)
		if err != nil {
			return false, err
		}

		return false, tryWrapDbError(
			r.client.
				Exec("DELETE FROM playlist_song_voters WHERE playlist_id = ? AND song_id = ?", playlistId, songId).
				Error,
		)
	}
}

func (r *Repository) AddSongToFavorites(songId, accountId uint) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.FavoriteSong)).
			Create(
				&models.FavoriteSong{
					AccountId: accountId,
					SongId:    songId,
				}).
			Error,
	)
}

func (r *Repository) RemoveSongFromFavorites(songId, accountId uint) error {
	return tryWrapDbError(
		r.client.
			Exec("DELETE FROM favorite_songs WHERE song_id = ? AND account_id = ?", songId, accountId).
			Error,
	)
}

func (r *Repository) GetFavoriteSongs(accountId, page uint) (models.List[models.Song], error) {
	gigaQuery := fmt.Sprintf(
		`SELECT songs.id, public_id, title, artist, thumbnail_url, real_duration, f.created_at
		FROM
			favorite_songs f JOIN songs
		ON
				songs.id = f.song_id
		WHERE f.account_id = ?
		ORDER BY f.created_at DESC
		LIMIT %d,%d;`,
		(page-1)*20, page*20,
	)

	rows, err := r.client.
		Raw(gigaQuery, accountId).
		Rows()
	if err != nil {
		return models.List[models.Song]{}, err
	}

	songs := make([]models.Song, 0)
	for rows.Next() {
		var song models.Song
		var addedAt time.Time
		err = rows.Scan(&song.Id, &song.PublicId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.RealDuration, &addedAt)
		if err != nil {
			continue
		}
		song.AddedAt = whenDidItHappen(addedAt)
		song.Favorite = true
		songs = append(songs, song)
	}
	_ = rows.Close()

	return models.NewList(songs, fmt.Sprint(page+1)), nil
}

func (r *Repository) IsSongFavorite(accountId, songId uint) error {
	var fav models.FavoriteSong
	err := tryWrapDbError(
		r.client.
			Model(new(models.FavoriteSong)).
			First(&fav, "account_id = ? AND song_id = ?", accountId, songId).
			Error,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPlaylistSongs(playlistId uint) (models.List[*models.Song], error) {
	gigaQuery := `SELECT songs.id, public_id, title, artist, thumbnail_url, real_duration, ps.created_at, ps.play_times, ps.votes
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
	songIds := make([]uint, 0)
	for rows.Next() {
		var song models.Song
		var addedAt time.Time
		err = rows.Scan(&song.Id, &song.PublicId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.RealDuration, &addedAt, &song.PlayTimes, &song.Votes)
		if err != nil {
			continue
		}
		song.AddedAt = addedAt.Format("2, January, 2006")
		songs = append(songs, &song)
		songIds = append(songIds, song.Id)
	}

	rows, err = r.client.
		Raw(`SELECT song_id FROM favorite_songs WHERE song_id IN ?`, songIds).
		Rows()
	if err != nil {
		return models.List[*models.Song]{}, err
	}

	songInFavorites := map[uint]bool{}
	for rows.Next() {
		var songId uint
		err = rows.Scan(&songId)
		if err != nil {
			continue
		}
		songInFavorites[songId] = true
	}

	for i := range songs {
		songs[i].Favorite = songInFavorites[songs[i].Id]
	}

	_ = rows.Close()

	return models.NewList(songs, ""), nil
}

func (r *Repository) MarkSongAsDownloaded(songPublicId string) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Song)).
			Where("public_id = ?", songPublicId).
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

func (r *Repository) AddAccountToPlaylist(plId, accountId uint, permissions models.PlaylistPermissions) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.PlaylistOwner)).
			Create(&models.PlaylistOwner{
				PlaylistId:  plId,
				AccountId:   accountId,
				Permissions: permissions,
			}).
			Error,
	)
}

func (r *Repository) RemoveAccountFromPlaylist(plId, accountId uint) error {
	return tryWrapDbError(
		r.client.
			Model(new(models.PlaylistOwner)).
			Delete(&models.PlaylistOwner{},
				"account_id = ? AND playlist_id = ?", accountId, plId,
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

func (r *Repository) GetPlaylistsForAccount(accountId uint) (models.List[models.Playlist], error) {
	var playlists []models.Playlist
	err := tryWrapDbError(
		r.client.
			Model(&models.Account{
				Id: accountId,
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

func (r *Repository) GetPlaylistsWithSongsForAccount(accountId uint) (models.List[models.Playlist], error) {
	var dbPlaylists []models.Playlist
	err := tryWrapDbError(
		r.client.
			Model(&models.Account{
				Id: accountId,
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

func (r *Repository) IncrementSongsCountForPlaylist(id uint) error {
	updateQuery := `UPDATE playlists
		SET songs_count = songs_count + 1
		WHERE
			id = ?;`

	err := tryWrapDbError(
		r.client.
			Exec(updateQuery, id).
			Error,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DecrementSongsCountForPlaylist(id uint) error {
	updateQuery := `UPDATE playlists
		SET songs_count = songs_count - 1
		WHERE
			id = ?;`

	err := tryWrapDbError(
		r.client.
			Exec(updateQuery, id).
			Error,
	)
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
			Limit(int(size)).
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

package mariadb

import (
	"dankmuzikk/app/models"
	"dankmuzikk/evy"
)

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.Account),
		new(models.Profile),
		new(models.Song),
		new(models.Playlist),
		new(models.PlaylistSong),
		new(models.PlaylistOwner),
		new(models.History),
		new(models.PlaylistSongVoter),
		new(evy.EventPayload),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"profiles", "songs", "playlists", "event_payloads",
	} {
		err = dbConn.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	return nil
}

func Migrate2() error {
	oldDb, err := dbConnector2()
	if err != nil {
		return err
	}

	newDb, err := dbConnector()
	if err != nil {
		return err
	}

	err = newDb.Debug().AutoMigrate(
		new(models.Account),
		new(models.Profile),
		new(models.Song),
		new(models.History),
		new(models.Playlist),
		new(models.PlaylistSong),
		new(models.PlaylistOwner),

		new(models.PlaylistSongVoter),

		new(evy.EventPayload),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"profiles", "songs", "playlists", "event_payloads",
	} {
		err = newDb.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	// users

	accountToProfileId := map[uint]uint{}
	profileToAccountId := map[uint]uint{}

	var accounts []models.Account
	rows, err := oldDb.Raw(`select id, email, is_o_auth, updated_at, created_at from accounts`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var a models.Account
		err = rows.Scan(&a.Id, &a.Email, &a.IsOAuth, &a.UpdatedAt, &a.CreatedAt)
		if err != nil {
			return err
		}
		accounts = append(accounts, a)
	}

	var profiles []models.Profile
	rows, err = oldDb.Raw(`select id, account_id, name, username, pfp_link, updated_at, created_at from profiles`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var a models.Profile
		err = rows.Scan(&a.Id, &a.AccountId, &a.Name, &a.Username, &a.PfpLink, &a.UpdatedAt, &a.CreatedAt)
		if err != nil {
			return err
		}

		accountToProfileId[a.AccountId] = a.Id
		profileToAccountId[a.Id] = a.AccountId

		profiles = append(profiles, a)
	}

	for _, a := range accounts {
		err = newDb.Exec(`insert into accounts (id, email, is_o_auth, updated_at, created_at) values (?, ?, ?, ?, ?);`,
			accountToProfileId[a.Id], a.Email, a.IsOAuth, a.UpdatedAt, a.CreatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	for _, a := range profiles {
		err = newDb.Exec(`insert into profiles (id, account_id, name, username, pfp_link, updated_at, created_at) values (?, ?, ?, ?, ?, ?, ?);`,
			profileToAccountId[a.Id], accountToProfileId[a.AccountId], a.Name, a.Username, a.PfpLink, a.UpdatedAt, a.CreatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	// history

	var historyItems []models.History
	rows, err = oldDb.Raw(`select id, song_id, profile_id, created_at from histories`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var h models.History
		err = rows.Scan(&h.Id, &h.SongId, &h.AccountId, &h.CreatedAt)
		if err != nil {
			return err
		}
		historyItems = append(historyItems, h)
	}

	for _, h := range historyItems {
		err = newDb.Exec(`insert into histories (id, song_id, account_id, created_at) values (?, ?, ?, ?);`,
			h.Id, h.SongId, h.AccountId, h.CreatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	// songs

	var songs []models.Song

	rows, err = oldDb.Raw(`select id, yt_id, title, artist, thumbnail_url, duration, fully_downloaded, created_at, updated_at from songs`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var h models.Song
		err = rows.Scan(&h.Id, &h.YtId, &h.Title, &h.Artist, &h.ThumbnailUrl, &h.Duration, &h.FullyDownloaded, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return err
		}
		songs = append(songs, h)
	}

	for _, h := range songs {
		err = newDb.Exec(`insert into songs (id, yt_id, title, artist, thumbnail_url, duration, fully_downloaded, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
			h.Id, h.YtId, h.Title, h.Artist, h.ThumbnailUrl, h.Duration, h.FullyDownloaded, h.CreatedAt, h.UpdatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	// playlists

	var playlists []models.Playlist

	rows, err = oldDb.Raw(`select id, public_id, title, songs_count, is_public, created_at, updated_at from playlists`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var h models.Playlist
		err = rows.Scan(&h.Id, &h.PublicId, &h.Title, &h.SongsCount, &h.IsPublic, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return err
		}
		playlists = append(playlists, h)
	}

	for _, h := range playlists {
		err = newDb.Exec(`insert into playlists (id, public_id, title, songs_count, is_public, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?);`,
			h.Id, h.PublicId, h.Title, h.SongsCount, h.IsPublic, h.CreatedAt, h.UpdatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	var playlistSongs []models.PlaylistSong

	rows, err = oldDb.Raw(`select playlist_id, song_id, votes, play_times, created_at, updated_at from playlist_songs`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var h models.PlaylistSong
		err = rows.Scan(&h.PlaylistId, &h.SongId, &h.Votes, &h.PlayTimes, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return err
		}
		playlistSongs = append(playlistSongs, h)
	}

	for _, h := range playlistSongs {
		err = newDb.Exec(`insert into playlist_songs (playlist_id, song_id, votes, play_times, created_at, updated_at) values (?, ?, ?, ?, ?, ?);`,
			h.PlaylistId, h.SongId, h.Votes, h.PlayTimes, h.CreatedAt, h.UpdatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	var playlistOwners []models.PlaylistOwner

	rows, err = oldDb.Raw(`select playlist_id, profile_id, permissions, created_at, updated_at from playlist_owners`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var h models.PlaylistOwner
		err = rows.Scan(&h.PlaylistId, &h.AccountId, &h.Permissions, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return err
		}
		playlistOwners = append(playlistOwners, h)
	}

	for _, h := range playlistOwners {
		err = newDb.Exec(`insert into playlist_owners (playlist_id, account_id, permissions, created_at, updated_at) values (?, ?, ?, ?, ?);`,
			h.PlaylistId, h.AccountId, h.Permissions, h.CreatedAt, h.UpdatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	var playlistVoters []models.PlaylistSongVoter

	rows, err = oldDb.Raw(`select playlist_id, song_id, profile_id, vote_up, created_at, updated_at from playlist_song_voters`).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var h models.PlaylistSongVoter
		err = rows.Scan(&h.PlaylistId, &h.SongId, &h.AccountId, &h.VoteUp, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return err
		}
		playlistVoters = append(playlistVoters, h)
	}

	for _, h := range playlistVoters {
		err = newDb.Exec(`insert into playlist_song_voters (playlist_id, song_id, account_id, vote_up, created_at, updated_at) values (?, ?, ?, ?, ?, ?);`,
			h.PlaylistId, h.SongId, h.AccountId, h.VoteUp, h.CreatedAt, h.UpdatedAt,
		).Error
		if err != nil {
			return err
		}
	}

	return nil
}

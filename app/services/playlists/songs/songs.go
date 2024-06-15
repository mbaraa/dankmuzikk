package songs

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/youtube/download"
	"dankmuzikk/services/youtube/search"
	"errors"
)

// Service represents songs in platlists management service,
// where it adds and deletes songs to and from playlists
type Service struct {
	playlistSongRepo   db.UnsafeCRUDRepo[models.PlaylistSong]
	playlistOwnerRepo  db.CRUDRepo[models.PlaylistOwner]
	songRepo           db.UnsafeCRUDRepo[models.Song]
	playlistRepo       db.UnsafeCRUDRepo[models.Playlist]
	playlistVotersRepo db.CRUDRepo[models.PlaylistSongVoter]
	downloadService    *download.Service
}

// New accepts repos lol, and returns a new instance to the songs playlists service.
func New(
	playlistSongRepo db.UnsafeCRUDRepo[models.PlaylistSong],
	playlistOwnerRepo db.CRUDRepo[models.PlaylistOwner],
	songRepo db.UnsafeCRUDRepo[models.Song],
	playlistRepo db.UnsafeCRUDRepo[models.Playlist],
	playlistVotersRepo db.CRUDRepo[models.PlaylistSongVoter],
	downloadService *download.Service,
) *Service {
	return &Service{
		playlistSongRepo:   playlistSongRepo,
		playlistOwnerRepo:  playlistOwnerRepo,
		songRepo:           songRepo,
		playlistRepo:       playlistRepo,
		playlistVotersRepo: playlistVotersRepo,
		downloadService:    downloadService,
	}
}

// GetSong returns a song with the provided youtube id, and an occurring error
func (s *Service) GetSong(songYtId string) (entities.Song, error) {
	song, err := s.songRepo.GetByConds("yt_id = ?", songYtId)
	if err != nil && errors.Is(err, db.ErrRecordNotFound) {
		res, err := (&search.ScraperSearch{}).Search(songYtId)
		if err != nil {
			return entities.Song{}, err
		}
		if len(res) == 0 {
			return entities.Song{}, errors.New("no songs were found suka")
		}
		for _, sng := range res {
			if sng.YtId == songYtId {
				ss := models.Song{
					YtId:            sng.YtId,
					Title:           sng.Title,
					Artist:          sng.Artist,
					ThumbnailUrl:    sng.ThumbnailUrl,
					Duration:        sng.Duration,
					FullyDownloaded: false,
				}
				err = s.songRepo.Add(&ss)
				log.Errorln(err)
				if len(song) == 0 {
					song = make([]models.Song, 1)
				}
				song[0] = ss
			}
		}
		err = s.downloadService.DownloadYoutubeSong(songYtId)
		if err != nil {
			return entities.Song{}, err
		}
	} else if err != nil {
		return entities.Song{}, err
	}
	if len(song) == 0 {
		return entities.Song{}, db.ErrRecordNotFound
	}

	return entities.Song{
		YtId:         song[0].YtId,
		Title:        song[0].Title,
		Artist:       song[0].Artist,
		ThumbnailUrl: song[0].ThumbnailUrl,
		Duration:     song[0].Duration,
	}, nil
}

// ToggleSongInPlaylist adds/removes a given song to/from the given playlist,
// checks if the actual song and playlist exist then adds/removes the song to/from the given playlist,
// and returns an occurring error.
func (s *Service) ToggleSongInPlaylist(songId, playlistPubId string, ownerId uint) (added bool, err error) {
	playlist, err := s.playlistRepo.GetByConds("public_id = ?", playlistPubId)
	if err != nil {
		return
	}
	_, err = s.playlistOwnerRepo.GetByConds("profile_id = ? AND playlist_id = ?", ownerId, playlist[0].Id)
	if err != nil {
		return
	}
	song, err := s.songRepo.GetByConds("yt_id = ?", songId)
	if err != nil {
		return
	}
	_, err = s.playlistSongRepo.GetByConds("playlist_id = ? AND song_id = ?", playlist[0].Id, song[0].Id)
	if errors.Is(err, db.ErrRecordNotFound) {
		err = s.playlistSongRepo.Add(&models.PlaylistSong{
			PlaylistId: playlist[0].Id,
			SongId:     song[0].Id,
			Votes:      1,
		})
		if err != nil {
			return
		}
		return true, s.downloadService.DownloadYoutubeSongQueue(songId)
	} else {
		return false, s.
			playlistSongRepo.
			Delete("playlist_id = ? AND song_id = ?", playlist[0].Id, song[0].Id)
	}
}

// IncrementSongPlays increases the song's play times in the given playlist.
// Checks for the song and playlist first, yada yada...
func (s *Service) IncrementSongPlays(songId, playlistPubId string, ownerId uint) error {
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
	err := s.songRepo.
		GetDB().
		Raw(gigaQuery, playlistPubId, songId, ownerId).
		Row().
		Scan(&playlistDbId, &songDbId)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE playlist_songs
		SET play_times = play_times + 1
		WHERE
			playlist_id = ? AND song_id = ?;`

	err = s.songRepo.
		GetDB().
		Exec(updateQuery, playlistDbId, songDbId).
		Error
	if err != nil {
		return err
	}

	return nil
}

// UpvoteSong increases the song's votes in the given playlist.
// Checks for the song and playlist first, yada yada...
func (s *Service) UpvoteSong(songId, playlistPubId string, ownerId uint) (int, error) {
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
	err := s.songRepo.
		GetDB().
		Raw(gigaQuery, playlistPubId, songId, ownerId).
		Row().
		Scan(&playlistDbId, &songDbId)
	if err != nil {
		return 0, err
	}

	voter, err := s.
		playlistVotersRepo.
		GetByConds("playlist_id = ? AND song_id = ? AND profile_id = ? AND vote_up = 1",
			playlistDbId,
			songDbId,
			ownerId,
		)
	if err == nil && len(voter) != 0 {
		return 0, playlists.ErrUserHasAlreadyVoted
	}

	updateQuery := `UPDATE playlist_songs
		SET votes = votes + 1
		WHERE
			playlist_id = ? AND song_id = ?;`

	err = s.songRepo.
		GetDB().
		Exec(updateQuery, playlistDbId, songDbId).
		Error
	if err != nil {
		return 0, err
	}

	ps, err := s.playlistSongRepo.GetByConds("playlist_id = ? AND song_id = ?", playlistDbId, songDbId)
	if err != nil {
		return 0, err
	}

	err = s.playlistVotersRepo.Add(&models.PlaylistSongVoter{
		PlaylistId: playlistDbId,
		SongId:     songDbId,
		ProfileId:  ownerId,
		VoteUp:     true,
	})
	log.Warningf("%+v\n%v\n", ps, err)
	if errors.Is(err, db.ErrRecordExists) {
		return ps[0].Votes, s.
			songRepo.
			GetDB().
			Exec("UPDATE playlist_song_voters SET vote_up = 1 WHERE playlist_id = ? AND song_id = ? AND profile_id = ?", playlistDbId, songDbId, ownerId).
			Error
	} else {
		return ps[0].Votes, err
	}
}

// DownvoteSong decreases the song's votes in the given playlist.
// Checks for the song and playlist first, yada yada...
func (s *Service) DownvoteSong(songId, playlistPubId string, ownerId uint) (int, error) {
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
	err := s.songRepo.
		GetDB().
		Raw(gigaQuery, playlistPubId, songId, ownerId).
		Row().
		Scan(&playlistDbId, &songDbId)
	if err != nil {
		return 0, err
	}

	voter, err := s.
		playlistVotersRepo.
		GetByConds("playlist_id = ? AND song_id = ? AND profile_id = ? AND vote_up = 0",
			playlistDbId,
			songDbId,
			ownerId,
		)
	if err == nil && len(voter) != 0 {
		return 0, playlists.ErrUserHasAlreadyVoted
	}
	log.Warningf("suka: %+v\n", voter)

	updateQuery := `UPDATE playlist_songs
		SET votes = votes - 1
		WHERE
			playlist_id = ? AND song_id = ?;`
	err = s.songRepo.
		GetDB().
		Exec(updateQuery, playlistDbId, songDbId).
		Error
	if err != nil {
		return 0, err
	}

	// remove song from playlist if votes < 0
	ps, err := s.playlistSongRepo.GetByConds("playlist_id = ? AND song_id = ?", playlistDbId, songDbId)
	if err != nil {
		return 0, err
	}
	if ps[0].Votes < 0 {
		return 0, s.playlistSongRepo.Delete("playlist_id = ? AND song_id = ?", playlistDbId, songDbId)
	}

	err = s.playlistVotersRepo.Add(&models.PlaylistSongVoter{
		PlaylistId: playlistDbId,
		SongId:     songDbId,
		ProfileId:  ownerId,
		VoteUp:     false,
	})
	log.Warningf("%+v\n%v\n", ps, err)
	if errors.Is(err, db.ErrRecordExists) {
		return ps[0].Votes, s.
			songRepo.
			GetDB().
			Exec("UPDATE playlist_song_voters SET vote_up = 0 WHERE playlist_id = ? AND song_id = ? AND profile_id = ?", playlistDbId, songDbId, ownerId).
			Error
	} else {
		return ps[0].Votes, err
	}
}

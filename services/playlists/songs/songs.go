package songs

import (
	"dankmuzikk/db"
	"dankmuzikk/models"
	"dankmuzikk/services/youtube/download"
	"errors"
)

// Service represents songs in platlists management service,
// where it adds and deletes songs to and from playlists
type Service struct {
	playlistSongRepo  db.UnsafeCRUDRepo[models.PlaylistSong]
	playlistOwnerRepo db.CRUDRepo[models.PlaylistOwner]
	songRepo          db.UnsafeCRUDRepo[models.Song]
	playlistRepo      db.UnsafeCRUDRepo[models.Playlist]
	downloadService   *download.Service
}

// New accepts repos lol, and returns a new instance to the songs playlists service.
func New(
	playlistSongRepo db.UnsafeCRUDRepo[models.PlaylistSong],
	playlistOwnerRepo db.CRUDRepo[models.PlaylistOwner],
	songRepo db.UnsafeCRUDRepo[models.Song],
	playlistRepo db.UnsafeCRUDRepo[models.Playlist],
	downloadService *download.Service,
) *Service {
	return &Service{
		playlistSongRepo:  playlistSongRepo,
		playlistOwnerRepo: playlistOwnerRepo,
		songRepo:          songRepo,
		playlistRepo:      playlistRepo,
		downloadService:   downloadService,
	}
}

// song-id vBHild0PiTE AND playlist-id 9d89cb71e00249cc99e1e5c54f3b0d5f

// ToggleSongInPlaylist adds/removes a given song to/from the given playlist,
// checks if the actual song and playlist exist then adds/removes the song to/from the given playlist,
// and returns an occurring error.
func (s *Service) ToggleSongInPlaylist(songId, playlistPubId string, ownerId uint) (added bool, err error) {
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
	err = s.songRepo.
		GetDB().
		Raw(gigaQuery, playlistPubId, songId, ownerId).
		Row().
		Scan(&playlistDbId, &songDbId)
	if err != nil {
		return false, err
	}

	_, err = s.playlistSongRepo.GetByConds("playlist_id = ? AND song_id = ?", playlistDbId, songDbId)
	if errors.Is(err, db.ErrRecordNotFound) {
		err = s.playlistSongRepo.Add(&models.PlaylistSong{
			PlaylistId: playlistDbId,
			SongId:     songDbId,
		})
		if err != nil {
			return
		}
		return true, s.downloadService.DownloadYoutubeSongQueue(songId)
	} else {
		return false, s.
			playlistSongRepo.
			Delete("playlist_id = ? AND song_id = ?", playlistDbId, songDbId)
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

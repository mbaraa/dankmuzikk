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
	var playlist models.Playlist
	err := s.
		songRepo.
		GetDB().
		Model(new(models.Playlist)).
		Select("id").
		Where("public_id = ?", playlistPubId).
		First(&playlist).
		Error
	if err != nil {
		return err
	}

	_, err = s.playlistOwnerRepo.GetByConds("profile_id = ? AND playlist_id = ?", ownerId, playlist.Id)
	if err != nil {
		return err
	}

	var song models.Song
	err = s.
		songRepo.
		GetDB().
		Model(new(models.Song)).
		Select("id").
		Where("yt_id = ?", songId).
		First(&song).
		Error
	if err != nil {
		return err
	}

	var ps models.PlaylistSong
	err = s.
		playlistSongRepo.
		GetDB().
		Model(new(models.PlaylistSong)).
		Select("play_times").
		Where("playlist_id = ? AND song_id = ?", playlist.Id, song.Id).
		First(&ps).
		Error
	if err != nil {
		return err
	}

	return s.
		playlistSongRepo.
		GetDB().
		Model(new(models.PlaylistSong)).
		Where("playlist_id = ? AND song_id = ?", playlist.Id, song.Id).
		Update("play_times", ps.PlayTimes+1).
		Error
}

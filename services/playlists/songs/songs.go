package songs

import (
	"dankmuzikk/db"
	"dankmuzikk/models"
)

// Service represents songs in platlists management service,
// where it adds and deletes songs to and from playlists
type Service struct {
	playlistSongRepo db.UnsafeCRUDRepo[models.PlaylistSong]
	songRepo         db.UnsafeCRUDRepo[models.Song]
	playlistRepo     db.UnsafeCRUDRepo[models.Playlist]
}

// New accepts repos lol, and returns a new instance to the songs playlists service.
func New(
	playlistSongRepo db.UnsafeCRUDRepo[models.PlaylistSong],
	songRepo db.UnsafeCRUDRepo[models.Song],
	playlistRepo db.UnsafeCRUDRepo[models.Playlist],
) *Service {
	return &Service{playlistSongRepo, songRepo, playlistRepo}
}

// AddSongToPlaylist adds a given song to the given playlist,
// checks if the actual song and playlist exist then adds the song to the given playlist,
// and returns an occurring error.
// TODO: check playlist's owner :)
func (s *Service) AddSongToPlaylist(songId, playlistPubId string) error {
	song, err := s.songRepo.GetByConds("yt_id = ?", songId)
	if err != nil {
		return err
	}
	playlist, err := s.playlistRepo.GetByConds("public_id = ?", playlistPubId)
	if err != nil {
		return err
	}

	return s.playlistSongRepo.Add(&models.PlaylistSong{
		PlaylistId: playlist[0].Id,
		SongId:     song[0].Id,
	})
}

// IncrementSongPlays increases the song's play times in the given playlist.
// Checks for the song and playlist first, yada yada...
// TODO: check playlist's owner :)
func (s *Service) IncrementSongPlays(songId, playlistPubId string) error {
	var song models.Song
	err := s.
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

	var playlist models.Playlist
	err = s.
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

// RemoveSongFromPlaylist removes a given song from the given playlist,
// checks if the actual song and playlist exist then removes the song to the given playlist,
// and returns an occurring error.
// TODO: check playlist's owner :)
func (s *Service) RemoveSongFromPlaylist(songId, playlistPubId string) error {
	song, err := s.songRepo.GetByConds("yt_id = ?", songId)
	if err != nil {
		return err
	}
	playlist, err := s.playlistRepo.GetByConds("public_id = ?", playlistPubId)
	if err != nil {
		return err
	}

	return s.
		playlistSongRepo.
		Delete("playlist_id = ? AND song_id = ?", playlist[0].Id, song[0].Id)
}

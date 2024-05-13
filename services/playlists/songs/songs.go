package songs

import (
	"dankmuzikk/db"
	"dankmuzikk/models"
)

// Service represents songs in platlists management service,
// where it adds and deletes songs to and from playlists
type Service struct {
	repo         db.CRUDRepo[models.PlaylistSong]
	songRepo     db.GetterRepo[models.Song]
	playlistRepo db.GetterRepo[models.Playlist]
}

// New accepts repos lol, and returns a new instance to the songs playlists service.
func New(
	repo db.CRUDRepo[models.PlaylistSong],
	songRepo db.GetterRepo[models.Song],
	playlistRepo db.GetterRepo[models.Playlist],
) *Service {
	return &Service{repo, songRepo, playlistRepo}
}

// AddSongToPlaylist adds a given song to the given playlist,
// checks if the actual song and playlist exist then adds the song to the given playlist,
// and returns an occurring error.
func (s *Service) AddSongToPlaylist(songId, playlistId uint) error {
	err := s.songRepo.Exists(songId)
	if err != nil {
		return err
	}
	err = s.playlistRepo.Exists(playlistId)
	if err != nil {
		return err
	}

	return s.repo.Add(&models.PlaylistSong{
		PlaylistId: playlistId,
		SongId:     songId,
	})
}

// RemoveSongFromPlaylist removes a given song from the given playlist,
// checks if the actual song and playlist exist then removes the song to the given playlist,
// and returns an occurring error.
func (s *Service) RemoveSongFromPlaylist(songId, playlistId uint) error {
	err := s.songRepo.Exists(songId)
	if err != nil {
		return err
	}
	err = s.playlistRepo.Exists(playlistId)
	if err != nil {
		return err
	}

	return s.
		repo.
		Delete("playlist_id = ? AND song_id = ?", playlistId, songId)
}
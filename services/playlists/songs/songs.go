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

	return s.repo.Add(&models.PlaylistSong{
		PlaylistId: playlist[0].Id,
		SongId:     song[0].Id,
	})
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
		repo.
		Delete("playlist_id = ? AND song_id = ?", playlist[0].Id, song[0].Id)
}

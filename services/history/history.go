package history

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"errors"
)

type Service struct {
	repo     db.UnsafeCRUDRepo[models.History]
	songRepo db.GetterRepo[models.Song]
}

func New(repo db.UnsafeCRUDRepo[models.History], songRepo db.GetterRepo[models.Song]) *Service {
	return &Service{
		repo:     repo,
		songRepo: songRepo,
	}
}

func (h *Service) AddSongToHistory(songYtId string, profileId uint) error {
	song, err := h.songRepo.GetByConds("yt_id = ?", songYtId)
	if err != nil {
		return err
	}

	return h.repo.Add(&models.History{
		ProfileId: profileId,
		SongId:    song[0].Id,
	})
}

func (h *Service) Get() ([]entities.Song, error) {
	return nil, errors.New("not implemented")
}

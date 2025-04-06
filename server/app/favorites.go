package app

import "dankmuzikk/app/models"

func (a *App) AddSongToFavorites(songPublicId string, accountId uint) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	return a.repo.AddSongToFavorites(song.Id, accountId)
}

func (a *App) RemoveSongFromFavorites(songPublicId string, accountId uint) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	return a.repo.RemoveSongFromFavorites(song.Id, accountId)
}

func (a *App) GetFavoriteSongs(page, accountId uint) (models.List[models.Song], error) {
	return a.repo.GetFavoriteSongs(accountId, page)
}

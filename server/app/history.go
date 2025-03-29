package app

import "dankmuzikk/app/models"

func (a *App) GetHistory(accountId, page uint) (models.List[models.Song], error) {
	songs, err := a.repo.GetHistory(accountId, page)
	if err != nil {
		return models.List[models.Song]{}, err
	}

	return songs, err
}

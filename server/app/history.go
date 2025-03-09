package app

import "dankmuzikk/app/models"

func (a *App) GetHistory(profileId, page uint) (models.List[models.Song], error) {
	songs, err := a.repo.GetHistory(profileId, page)
	if err != nil {
		return models.List[models.Song]{}, err
	}

	return songs, err
}

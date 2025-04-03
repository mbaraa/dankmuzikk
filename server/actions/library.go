package actions

type AddSongToFavoritesParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) AddSongToFavorites(params AddSongToFavoritesParams) error {
	return a.app.AddSongToFavorites(params.SongPublicId, params.Account.Id)
}

type RemoveSongFromFavoritesParams struct {
	ActionContext `json:"-"`
	SongPublicId  string `json:"song_public_id"`
}

func (a *Actions) RemoveSongFromFavorites(params RemoveSongFromFavoritesParams) error {
	return a.app.RemoveSongFromFavorites(params.SongPublicId, params.Account.Id)
}

type GetFavoriteSongsParams struct {
	ActionContext `json:"-"`
	PageIndex     uint `json:"page_index"`
}

type GetFavoriteSongsPayload struct {
	Songs []Song `json:"songs"`
}

func (a *Actions) GetFavoriteSongs(params GetFavoriteSongsParams) (GetFavoriteSongsPayload, error) {
	favoriteSongs, err := a.app.GetFavoriteSongs(params.PageIndex, params.Account.Id)
	if err != nil {
		return GetFavoriteSongsPayload{}, err
	}

	outSongs := make([]Song, 0, favoriteSongs.Size)
	for _, song := range favoriteSongs.Items {
		outSongs = append(outSongs, mapModelToActionsSong(song))
	}

	return GetFavoriteSongsPayload{
		Songs: outSongs,
	}, nil
}

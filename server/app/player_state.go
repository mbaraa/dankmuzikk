package app

import (
	"dankmuzikk/app/models"
	"math"
	"math/rand/v2"
	"slices"
)

func (a *App) OverrideSongsQueue(accountId uint, initialSongPublicIds ...string) error {
	songs, err := a.repo.GetSongsByPublicIds(initialSongPublicIds)
	if err != nil {
		return err
	}

	songIds := make([]uint, 0, len(songs))
	for _, song := range songs {
		songIds = append(songIds, song.Id)
	}

	err = a.playerCache.ClearQueue(accountId)
	if err != nil {
		return err
	}

	err = a.playerCache.SetCurrentPlayingPlaylistInQueue(accountId, 0)
	if err != nil {
		return err
	}

	err = a.playerCache.SetCurrentPlayingSongInedxInQueue(accountId, 0)
	if err != nil {
		return err
	}

	err = a.playerCache.SetCurrentPlayingSongInedxInShuffledQueue(accountId, 0)
	if err != nil {
		return err
	}

	return a.playerCache.CreateSongsQueue(accountId, songIds...)
}

func (a *App) CreateSongsQueue(accountId uint, initialSongPublicIds ...string) error {
	songs, err := a.repo.GetSongsByPublicIds(initialSongPublicIds)
	if err != nil {
		return err
	}

	songsIds := make([]uint, 0, len(songs))
	for _, song := range songs {
		songsIds = append(songsIds, song.Id)
	}

	shuffled, _ := a.playerCache.GetShuffled(accountId)
	if shuffled {
		return a.playerCache.CreateSongsShuffledQueue(accountId, songsIds...)
	}

	return a.playerCache.CreateSongsQueue(accountId, songsIds...)
}

func (a *App) GetPlayerState(accountId uint) (models.PlayerState, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId)
	songs, _ := a.getSongsFromQueue(accountId, shuffled)
	currentSongIndex, _ := a.getCurrentPlayingSongIndex(accountId, shuffled)
	currentPlaylistId, _ := a.playerCache.GetCurrentPlayingPlaylistInQueue(accountId)
	loopMode, _ := a.playerCache.GetLoopMode(accountId)

	return models.PlayerState{
		Shuffled:          shuffled,
		CurrentSongIndex:  currentSongIndex,
		LoopMode:          loopMode,
		CurrentPlaylistId: currentPlaylistId,
		Songs:             songs,
	}, nil
}

func (a *App) AddSongToQueue(accountId uint, songPublicId string) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	shuffled, _ := a.playerCache.GetShuffled(accountId)
	if shuffled {
		// TODO: maybe move to an event?
		_ = a.playerCache.AddSongToShuffledQueue(accountId, song.Id)
	}

	return a.playerCache.AddSongToQueue(accountId, song.Id)
}

func (a *App) AddSongToQueueAfterCurrentSong(accountId uint, songPublicId string) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId)
	currentSongIndex, err := a.getCurrentPlayingSongIndex(accountId, shuffled)
	if err != nil {
		return a.CreateSongsQueue(accountId, songPublicId)
	}

	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	if shuffled {
		// TODO: maybe move to an event?
		_ = a.playerCache.AddSongToShuffledQueueAfterIndex(accountId, song.Id, currentSongIndex)
	}

	return a.playerCache.AddSongToQueueAfterIndex(accountId, song.Id, currentSongIndex)
}

func (a *App) AddPlaylistToQueue(accountId uint, playlistPublicId string) error {
	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	songs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}

	shuffled, _ := a.playerCache.GetShuffled(accountId)
	if shuffled {
		// TODO: maybe move to an event?
		for _, song := range songs.Items {
			_ = a.playerCache.AddSongToShuffledQueue(accountId, song.Id)
		}
	}

	for _, song := range songs.Items {
		_ = a.playerCache.AddSongToQueue(accountId, song.Id)
	}

	return nil
}

func (a *App) AddPlaylistToQueueAfterCurrentSong(accountId uint, playlistPublicId string) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId)
	currentSongIndex, err := a.getCurrentPlayingSongIndex(accountId, shuffled)
	if err != nil {
		return a.CreateSongsQueue(accountId, playlistPublicId)
	}

	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	songs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}
	slices.Reverse(songs.Items)

	if shuffled {
		// TODO: maybe move to an event?
		for _, song := range songs.Items {
			_ = a.playerCache.AddSongToShuffledQueueAfterIndex(accountId, song.Id, currentSongIndex)
		}
	}

	for _, song := range songs.Items {
		_ = a.playerCache.AddSongToQueueAfterIndex(accountId, song.Id, currentSongIndex)
	}

	return nil
}

func (a *App) RemoveSongFromQueue(songIndex int, accountId uint) error {
	return a.playerCache.RemoveSongFromQueue(songIndex, accountId)
}

func (a *App) SetCurrentPlayingSongInedxInQueue(accountId uint, songIndex int) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId)
	return a.setCurrentPlayingSongIndex(accountId, shuffled, songIndex)
}

func (a *App) ClearQueue(accountId uint) error {
	return a.playerCache.ClearQueue(accountId)
}

func (a *App) SetShuffledOn(accountId uint) error {
	songsInQueue, err := a.playerCache.GetSongsQueue(accountId)
	if err != nil {
		return err
	}

	rand.Shuffle(len(songsInQueue), func(i, j int) {
		songsInQueue[i], songsInQueue[j] = songsInQueue[j], songsInQueue[i]
	})

	err = a.playerCache.ClearShuffledQueue(accountId)
	if err != nil {
		return err
	}

	err = a.playerCache.CreateSongsShuffledQueue(accountId, songsInQueue...)
	if err != nil {
		return err
	}

	return a.playerCache.SetShuffled(accountId, true)
}

func (a *App) SetShuffledOff(accountId uint) error {
	_ = a.playerCache.ClearShuffledQueue(accountId)

	return a.playerCache.SetShuffled(accountId, false)
}

func (a *App) SetLoopMode(accountId uint, loopMode models.PlayerLoopMode) error {
	return a.playerCache.SetLoopMode(accountId, loopMode)
}

func (a *App) PlayPlaylist(accountId uint, playlistPublicId string) error {
	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	player, err := a.GetPlayerState(accountId)
	if err != nil {
		return err
	}

	_ = a.setCurrentPlayingSongIndex(accountId, player.Shuffled, 0)

	if player.CurrentPlaylistId == playlist.Id {
		return nil
	}

	playlistSongs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}

	songIds := make([]uint, 0, playlistSongs.Size)
	for _, song := range playlistSongs.Items {
		songIds = append(songIds, song.Id)
	}

	if player.Shuffled {
		err := a.playerCache.ClearShuffledQueue(accountId)
		if err != nil {
			return err
		}
		rand.Shuffle(len(songIds), func(i, j int) {
			songIds[i], songIds[j] = songIds[j], songIds[i]
		})
		err = a.playerCache.CreateSongsShuffledQueue(accountId, songIds...)
		if err != nil {
			return err
		}
	} else {
		err := a.playerCache.ClearQueue(accountId)
		if err != nil {
			return err
		}
		err = a.playerCache.CreateSongsQueue(accountId, songIds...)
		if err != nil {
			return err
		}
	}

	return a.playerCache.SetCurrentPlayingPlaylistInQueue(accountId, playlist.Id)
}

func (a *App) PlaySongFromPlaylist(accountId uint, songPublicId, playlistPublicId string) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	player, err := a.GetPlayerState(accountId)
	if err != nil {
		return err
	}

	songIndex := 0
	for i := len(player.Songs) - 1; i >= 0; i-- {
		if player.Songs[i].Id == song.Id {
			songIndex = i
			break
		}
	}

	_ = a.setCurrentPlayingSongIndex(accountId, player.Shuffled, songIndex)

	if player.CurrentPlaylistId == playlist.Id {
		return nil
	}

	playlistSongs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}

	songIds := make([]uint, 0, playlistSongs.Size)
	for _, song := range playlistSongs.Items {
		songIds = append(songIds, song.Id)
	}

	if player.Shuffled {
		err := a.playerCache.ClearShuffledQueue(accountId)
		if err != nil {
			return err
		}
		rand.Shuffle(len(songIds), func(i, j int) {
			songIds[i], songIds[j] = songIds[j], songIds[i]
		})
		err = a.playerCache.CreateSongsShuffledQueue(accountId, songIds...)
		if err != nil {
			return err
		}
	} else {
		err := a.playerCache.ClearQueue(accountId)
		if err != nil {
			return err
		}
		err = a.playerCache.CreateSongsQueue(accountId, songIds...)
		if err != nil {
			return err
		}
	}

	return a.playerCache.SetCurrentPlayingPlaylistInQueue(accountId, playlist.Id)
}

func (a *App) PlaySongFromFavorites(accountId uint, songPublicId string) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	favoriteSongs := make([]models.Song, 0)

	for page := 0; ; page++ {
		fs, err := a.repo.GetFavoriteSongs(accountId, uint(page))
		if err != nil {
			break
		}

		favoriteSongs = append(favoriteSongs, fs.Items...)
	}

	if len(favoriteSongs) == 0 {
		return &ErrNotFound{
			ResourceName: "favoriteSongs",
		}
	}

	songPublicIds := make([]string, 0, len(favoriteSongs))
	for _, song := range favoriteSongs {
		songPublicIds = append(songPublicIds, song.PublicId)
	}

	songIndex := 0
	for i := len(songPublicIds) - 1; i >= 0; i-- {
		if songPublicIds[i] == song.PublicId {
			songIndex = i
			break
		}
	}

	err = a.OverrideSongsQueue(accountId, songPublicIds...)
	if err != nil {
		return err
	}

	return a.setCurrentPlayingSongIndex(accountId, false, songIndex)
}

type GetNextPlayingSongResult struct {
	Song                    models.Song
	CurrentPlayingSongIndex int
	EndOfQueue              bool
}

func (a *App) GetNextPlayingSong(accountId uint) (GetNextPlayingSongResult, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId)
	currentSongIndex, _ := a.getCurrentPlayingSongIndex(accountId, shuffled)
	queueLen, _ := a.getQueueLength(accountId, shuffled)
	loopMode, _ := a.playerCache.GetLoopMode(accountId)

	endOfQueue := false
	switch loopMode {
	case models.LoopOffMode:
		if currentSongIndex+1 < int(queueLen) {
			currentSongIndex++
		} else {
			endOfQueue = true
		}
	case models.LoopOnceMode:
		endOfQueue = true
		break
	case models.LoopAllMode:
		currentSongIndex = (currentSongIndex + 1) % int(queueLen)
	}
	_ = a.setCurrentPlayingSongIndex(accountId, shuffled, currentSongIndex)

	songId, err := a.getSongAtIndexFromQueue(accountId, shuffled, currentSongIndex)
	if err != nil {
		return GetNextPlayingSongResult{}, err
	}

	song, err := a.repo.GetSong(songId)
	if err != nil {
		return GetNextPlayingSongResult{}, err
	}

	return GetNextPlayingSongResult{
		Song:                    song,
		CurrentPlayingSongIndex: currentSongIndex,
		EndOfQueue:              endOfQueue,
	}, nil
}

type GetPreviousPlayingSongResult struct {
	Song                    models.Song
	CurrentPlayingSongIndex int
	EndOfQueue              bool
}

func (a *App) GetPreviousPlayingSong(accountId uint) (GetPreviousPlayingSongResult, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId)
	currentSongIndex, _ := a.getCurrentPlayingSongIndex(accountId, shuffled)
	queueLen, _ := a.getQueueLength(accountId, shuffled)
	loopMode, _ := a.playerCache.GetLoopMode(accountId)

	endOfQueue := false
	switch loopMode {
	case models.LoopOffMode:
		if currentSongIndex > 0 {
			currentSongIndex--
		} else {
			endOfQueue = true
		}
	case models.LoopOnceMode:
		endOfQueue = true
		break
	case models.LoopAllMode:
		currentSongIndex = int(math.Abs(float64((currentSongIndex - 1) % int(queueLen))))
	}
	_ = a.setCurrentPlayingSongIndex(accountId, shuffled, currentSongIndex)

	songId, err := a.getSongAtIndexFromQueue(accountId, shuffled, currentSongIndex)
	if err != nil {
		return GetPreviousPlayingSongResult{}, err
	}

	song, err := a.repo.GetSong(songId)
	if err != nil {
		return GetPreviousPlayingSongResult{}, err
	}

	return GetPreviousPlayingSongResult{
		Song:                    song,
		CurrentPlayingSongIndex: currentSongIndex,
		EndOfQueue:              endOfQueue,
	}, nil
}

func (a *App) getSongAtIndexFromQueue(accountId uint, shuffled bool, index int) (uint, error) {
	if shuffled {
		return a.playerCache.GetSongIdAtIndexFromShuffledQueue(accountId, index)
	}

	return a.playerCache.GetSongIdAtIndexFromQueue(accountId, index)
}

func (a *App) getQueueLength(accountId uint, shuffled bool) (uint, error) {
	if shuffled {
		return a.playerCache.GetShuffledQueueLength(accountId)
	}

	return a.playerCache.GetQueueLength(accountId)
}

func (a *App) getCurrentPlayingSongIndex(accountId uint, shuffled bool) (int, error) {
	if shuffled {
		return a.playerCache.GetCurrentPlayingSongIndexInShuffledQueue(accountId)
	}

	return a.playerCache.GetCurrentPlayingSongIndexInQueue(accountId)
}

func (a *App) setCurrentPlayingSongIndex(accountId uint, shuffled bool, songIndex int) error {
	if shuffled {
		_ = a.playerCache.SetCurrentPlayingSongInedxInShuffledQueue(accountId, songIndex)
	}

	return a.playerCache.SetCurrentPlayingSongInedxInQueue(accountId, songIndex)
}

func (a *App) getSongsFromQueue(accountId uint, shuffled bool) ([]models.Song, error) {
	var songIds []uint
	if shuffled {
		songIds, _ = a.playerCache.GetSongsShuffledQueue(accountId)
	} else {
		songIds, _ = a.playerCache.GetSongsQueue(accountId)
	}

	songIdsOrdered := make(map[uint][]int)
	for idx, songId := range songIds {
		songIdsOrdered[songId] = append(songIdsOrdered[songId], idx)
	}

	// error is ignored because a player's state is allowed to have an empty queue.
	songs, _ := a.repo.GetSongsByIds(songIds)

	songsOrdered := make([]models.Song, len(songIds))
	for _, song := range songs {
		for _, songIdx := range songIdsOrdered[song.Id] {
			if songIdx >= len(songsOrdered) {
				continue
			}
			songsOrdered[songIdx] = song
		}
	}

	return songsOrdered, nil
}

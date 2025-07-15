package app

import (
	"dankmuzikk/app/models"
	"math"
	"math/rand/v2"
	"slices"
)

func (a *App) CreateSongsQueue(accountId uint, clientHash string, initialSongPublicIds []string) error {
	var songIds []uint
	if len(initialSongPublicIds) > 0 {
		songs, err := a.repo.GetSongsByPublicIds(initialSongPublicIds)
		if err != nil {
			return err
		}

		songsUnordered := make([]models.Song, 0, len(songs))
		for _, song := range songs {
			songsUnordered = append(songsUnordered, song)
		}

		songIdsOrdered := make(map[string][]int)
		for idx, song := range initialSongPublicIds {
			songIdsOrdered[song] = append(songIdsOrdered[song], idx)
		}

		songIds = make([]uint, len(initialSongPublicIds))
		for _, song := range songsUnordered {
			for _, songIdx := range songIdsOrdered[song.PublicId] {
				if songIdx >= len(songIds) {
					continue
				}
				songIds[songIdx] = song.Id
			}
		}
	}

	err := a.ClearQueue(accountId, clientHash)
	if err != nil {
		return err
	}

	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		rand.Shuffle(len(songIds), func(i, j int) {
			songIds[i], songIds[j] = songIds[j], songIds[i]
		})

		return a.playerCache.CreateSongsShuffledQueue(accountId, clientHash, songIds...)
	}

	return a.playerCache.CreateSongsQueue(accountId, clientHash, songIds...)
}

func (a *App) GetPlayerState(accountId uint, clientHash string) (models.PlayerState, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	songs, _ := a.getSongsFromQueue(accountId, clientHash)
	currentSongIndex, _ := a.getCurrentPlayingSongIndex(accountId, clientHash)
	currentPlaylistId, _ := a.playerCache.GetCurrentPlayingPlaylistInQueue(accountId, clientHash)
	loopMode, _ := a.playerCache.GetLoopMode(accountId, clientHash)

	return models.PlayerState{
		Shuffled:          shuffled,
		CurrentSongIndex:  currentSongIndex,
		LoopMode:          loopMode,
		CurrentPlaylistId: currentPlaylistId,
		Songs:             songs,
	}, nil
}

func (a *App) AddSongToQueue(accountId uint, clientHash, songPublicId string) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		_ = a.playerCache.AddSongToShuffledQueue(accountId, clientHash, song.Id)
	}

	return a.playerCache.AddSongToQueue(accountId, clientHash, song.Id)
}

func (a *App) AddSongToQueueAfterCurrentSong(accountId uint, clientHash, songPublicId string) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	currentSongIndex, err := a.getCurrentPlayingSongIndex(accountId, clientHash)
	if err != nil {
		return a.CreateSongsQueue(accountId, clientHash, []string{songPublicId})
	}

	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	if shuffled {
		_ = a.playerCache.AddSongToShuffledQueueAfterIndex(accountId, clientHash, song.Id, currentSongIndex)
	}

	return a.playerCache.AddSongToQueueAfterIndex(accountId, clientHash, song.Id, currentSongIndex)
}

func (a *App) AddPlaylistToQueue(accountId uint, clientHash, playlistPublicId string) error {
	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	songs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}

	songIds := make([]uint, 0, songs.Size)
	for _, song := range songs.Items {
		songIds = append(songIds, song.Id)
	}

	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		rand.Shuffle(len(songIds), func(i, j int) {
			songIds[i], songIds[j] = songIds[j], songIds[i]
		})

		return a.playerCache.AddSongsToShuffledQueue(accountId, clientHash, songIds...)
	}

	return a.playerCache.AddSongsToQueue(accountId, clientHash, songIds...)
}

func (a *App) AddPlaylistToQueueAfterCurrentSong(accountId uint, clientHash, playlistPublicId string) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	currentSongIndex, err := a.getCurrentPlayingSongIndex(accountId, clientHash)
	if err != nil {
		_ = a.CreateSongsQueue(accountId, clientHash, nil)
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
		// TODO: pipeline this...
		for _, song := range songs.Items {
			_ = a.playerCache.AddSongToShuffledQueueAfterIndex(accountId, clientHash, song.Id, currentSongIndex)
		}
	}

	// TODO: pipeline this...
	for _, song := range songs.Items {
		_ = a.playerCache.AddSongToQueueAfterIndex(accountId, clientHash, song.Id, currentSongIndex)
	}

	return nil
}

func (a *App) RemoveSongFromQueue(accountId uint, clientHash string, songIndex int) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		return a.playerCache.RemoveSongFromShuffledQueue(accountId, clientHash, songIndex)
	}
	return a.playerCache.RemoveSongFromQueue(accountId, clientHash, songIndex)
}

func (a *App) ClearQueue(accountId uint, clientHash string) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	_ = a.playerCache.SetCurrentPlayingPlaylistInQueue(accountId, clientHash, 0)
	if shuffled {
		_ = a.playerCache.SetCurrentPlayingSongIndexInShuffledQueue(accountId, clientHash, 0)
		return a.playerCache.ClearShuffledQueue(accountId, clientHash)
	}

	_ = a.playerCache.SetCurrentPlayingSongIndexInQueue(accountId, clientHash, 0)
	return a.playerCache.ClearQueue(accountId, clientHash)
}

func (a *App) SetShuffledOn(accountId uint, clientHash string) error {
	songsInQueue, err := a.playerCache.GetSongsQueue(accountId, clientHash)
	if err != nil {
		return err
	}

	rand.Shuffle(len(songsInQueue), func(i, j int) {
		songsInQueue[i], songsInQueue[j] = songsInQueue[j], songsInQueue[i]
	})

	err = a.playerCache.ClearShuffledQueue(accountId, clientHash)
	if err != nil {
		return err
	}

	err = a.playerCache.CreateSongsShuffledQueue(accountId, clientHash, songsInQueue...)
	if err != nil {
		return err
	}

	return a.playerCache.SetShuffled(accountId, clientHash, true)
}

func (a *App) SetShuffledOff(accountId uint, clientHash string) error {
	err := a.playerCache.ClearShuffledQueue(accountId, clientHash)
	if err != nil {
		return err
	}

	return a.playerCache.SetShuffled(accountId, clientHash, false)
}

func (a *App) SetLoopMode(accountId uint, clientHash string, loopMode models.PlayerLoopMode) error {
	return a.playerCache.SetLoopMode(accountId, clientHash, loopMode)
}

func (a *App) PlayPlaylist(accountId uint, clientHash, playlistPublicId string) error {
	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	player, err := a.GetPlayerState(accountId, clientHash)
	if err != nil {
		return err
	}

	err = a.setCurrentPlayingSongIndex(accountId, clientHash, 0)
	if err != nil {
		return err
	}

	if player.CurrentPlaylistId == playlist.Id {
		return nil
	}

	playlistSongs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}

	songIds := make([]string, 0, playlistSongs.Size)
	for _, song := range playlistSongs.Items {
		songIds = append(songIds, song.PublicId)
	}

	err = a.CreateSongsQueue(accountId, clientHash, songIds)
	if err != nil {
		return err
	}

	return a.playerCache.SetCurrentPlayingPlaylistInQueue(accountId, clientHash, playlist.Id)
}

func (a *App) PlaySongFromPlaylist(accountId uint, clientHash, songPublicId, playlistPublicId string) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	playlist, err := a.repo.GetPlaylistByPublicId(playlistPublicId)
	if err != nil {
		return err
	}

	playlistSongs, err := a.repo.GetPlaylistSongs(playlist.Id)
	if err != nil {
		return err
	}

	currentPlaylistId, _ := a.playerCache.GetCurrentPlayingPlaylistInQueue(accountId, clientHash)
	if currentPlaylistId != playlist.Id {
		songIds := make([]string, 0, playlistSongs.Size)
		for _, song := range playlistSongs.Items {
			songIds = append(songIds, song.PublicId)
		}

		err = a.CreateSongsQueue(accountId, clientHash, songIds)
		if err != nil {
			return err
		}

		err = a.playerCache.SetCurrentPlayingPlaylistInQueue(accountId, clientHash, playlist.Id)
		if err != nil {
			return err
		}
	}

	songIndex := 0
	for i := playlistSongs.Size - 1; i >= 0; i-- {
		if playlistSongs.Items[i].Id == song.Id {
			songIndex = i
			break
		}
	}

	return a.setCurrentPlayingSongIndex(accountId, clientHash, songIndex)
}

func (a *App) PlaySongFromFavorites(accountId uint, clientHash, songPublicId string) error {
	song, err := a.repo.GetSongByPublicId(songPublicId)
	if err != nil {
		return err
	}

	favoriteSongs := make([]models.Song, 0)
	for page := 1; ; page++ {
		fs, err := a.repo.GetFavoriteSongs(accountId, uint(page))
		if err != nil || fs.Size == 0 {
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

	err = a.CreateSongsQueue(accountId, clientHash, songPublicIds)
	if err != nil {
		return err
	}

	return a.setCurrentPlayingSongIndex(accountId, clientHash, songIndex)
}

func (a *App) PlaySongFromQueue(accountId uint, clientHash, songPublicId string) error {
	songs, err := a.getSongsFromQueue(accountId, clientHash)
	if err != nil {
		return err
	}

	songIndex := 0
	for i := len(songs) - 1; i >= 0; i-- {
		if songs[i].PublicId == songPublicId {
			songIndex = i
			break
		}
	}

	return a.setCurrentPlayingSongIndex(accountId, clientHash, songIndex)
}

type GetNextPlayingSongResult struct {
	Song                    models.Song
	CurrentPlayingSongIndex int
	EndOfQueue              bool
}

func (a *App) GetNextPlayingSong(accountId uint, clientHash string) (GetNextPlayingSongResult, error) {
	currentSongIndex, _ := a.getCurrentPlayingSongIndex(accountId, clientHash)
	queueLen, _ := a.getQueueLength(accountId, clientHash)
	loopMode, _ := a.playerCache.GetLoopMode(accountId, clientHash)

	endOfQueue := false
	switch loopMode {
	case models.LoopOffMode:
		if currentSongIndex+1 < int(queueLen) {
			currentSongIndex++
		} else {
			endOfQueue = true
		}
	case models.LoopOnceMode:
		break
	case models.LoopAllMode:
		currentSongIndex = (currentSongIndex + 1) % int(queueLen)
	}
	_ = a.setCurrentPlayingSongIndex(accountId, clientHash, currentSongIndex)

	songId, err := a.getSongAtIndexFromQueue(accountId, clientHash, currentSongIndex)
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

func (a *App) GetPreviousPlayingSong(accountId uint, clientHash string) (GetPreviousPlayingSongResult, error) {
	currentSongIndex, _ := a.getCurrentPlayingSongIndex(accountId, clientHash)
	queueLen, _ := a.getQueueLength(accountId, clientHash)
	loopMode, _ := a.playerCache.GetLoopMode(accountId, clientHash)

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
	case models.LoopAllMode:
		currentSongIndex = int(math.Abs(float64((currentSongIndex - 1) % int(queueLen))))
	}
	_ = a.setCurrentPlayingSongIndex(accountId, clientHash, currentSongIndex)

	songId, err := a.getSongAtIndexFromQueue(accountId, clientHash, currentSongIndex)
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

func (a *App) getSongAtIndexFromQueue(accountId uint, clientHash string, index int) (uint, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		return a.playerCache.GetSongIdAtIndexFromShuffledQueue(accountId, clientHash, index)
	}

	return a.playerCache.GetSongIdAtIndexFromQueue(accountId, clientHash, index)
}

func (a *App) getQueueLength(accountId uint, clientHash string) (uint, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		return a.playerCache.GetShuffledQueueLength(accountId, clientHash)
	}

	return a.playerCache.GetQueueLength(accountId, clientHash)
}

func (a *App) GetCurrentPlayingSong(accountId uint, clientHash string) (models.Song, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	index, err := a.getCurrentPlayingSongIndex(accountId, clientHash)
	if err != nil {
		return models.Song{}, err
	}

	var song models.Song
	if shuffled {
		songId, err := a.playerCache.GetSongIdAtIndexFromShuffledQueue(accountId, clientHash, index)
		if err != nil {
			return models.Song{}, err
		}
		song, err = a.repo.GetSong(songId)
		if err != nil {
			return models.Song{}, err
		}
	} else {
		songId, err := a.playerCache.GetSongIdAtIndexFromQueue(accountId, clientHash, index)
		if err != nil {
			return models.Song{}, err
		}
		song, err = a.repo.GetSong(songId)
		if err != nil {
			return models.Song{}, err
		}
	}

	return song, nil
}

func (a *App) getCurrentPlayingSongIndex(accountId uint, clientHash string) (int, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		return a.playerCache.GetCurrentPlayingSongIndexInShuffledQueue(accountId, clientHash)
	}

	return a.playerCache.GetCurrentPlayingSongIndexInQueue(accountId, clientHash)
}

func (a *App) setCurrentPlayingSongIndex(accountId uint, clientHash string, songIndex int) error {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	if shuffled {
		return a.playerCache.SetCurrentPlayingSongIndexInShuffledQueue(accountId, clientHash, songIndex)
	}

	return a.playerCache.SetCurrentPlayingSongIndexInQueue(accountId, clientHash, songIndex)
}

func (a *App) getSongsFromQueue(accountId uint, clientHash string) ([]models.Song, error) {
	shuffled, _ := a.playerCache.GetShuffled(accountId, clientHash)
	var songIds []uint
	var err error
	if shuffled {
		songIds, err = a.playerCache.GetSongsShuffledQueue(accountId, clientHash)
	} else {
		songIds, err = a.playerCache.GetSongsQueue(accountId, clientHash)
	}
	if err != nil {
		return nil, err
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

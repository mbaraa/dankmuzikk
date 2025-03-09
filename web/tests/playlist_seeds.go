package tests

import (
	"dankmuzikk-web/models"
	"dankmuzikk-web/services/nanoid"
	"time"
)

var playlists = []models.Playlist{
	{Title: "Favorites"},
	{Title: "Summer Vibes"},
	{Title: "Hip Hop Hits"},
	{Title: "Pop Classics"},
	{Title: "Empowering Anthems"},
	{Title: "Chill Out"},
	{Title: "Throwback Jams"},
	{Title: "Party Starters"},
	{Title: "Breakup Ballads"},
	{Title: "Road Trip Tunes"},
}

func initPlaylists() {
	for i := range playlists {
		playlists[i].Songs = RandomSongs(random.Intn(6))
		playlists[i].SongsCount = len(playlists[i].Songs)
		playlists[i].PublicId = nanoid.Generate()
		playlists[i].IsPublic = random.Int()%2 == 0
	}
}

func RandomPlaylist() models.Playlist {
	return playlists[random.Intn(len(playlists))]
}

func RandomPlaylists(amount int) []models.Playlist {
	randPlaylist := make([]models.Playlist, amount)
	for i := 0; i < amount; i++ {
		randPlaylist[i] = RandomPlaylist()
		random.Seed(time.Now().UnixMicro())
	}
	return randPlaylist
}

func Playlists() []models.Playlist {
	return playlists
}

package tests

import (
	"dankmuzikk-web/models"
	"dankmuzikk-web/services/nanoid"
	"time"
)

var songs = []models.Song{
	{Title: "Born to Die", Artist: "Lana Del Rey", Duration: "4:45"},
	{Title: "Video Games", Artist: "Lana Del Rey", Duration: "4:02"},
	{Title: "Summertime Sadness", Artist: "Lana Del Rey", Duration: "4:10"},
	{Title: "Blue Jeans", Artist: "Lana Del Rey", Duration: "3:38"},
	{Title: "Lose Yourself", Artist: "Eminem", Duration: "4:23"},
	{Title: "Not Afraid", Artist: "Eminem", Duration: "4:11"},
	{Title: "Stan", Artist: "Eminem", Duration: "5:46"},
	{Title: "Love the Way You Lie", Artist: "Eminem", Duration: "3:35"},
	{Title: "Rolling in the Deep", Artist: "Adele", Duration: "3:41"},
	{Title: "Someone Like You", Artist: "Adele", Duration: "4:11"},
	{Title: "Set Fire to the Rain", Artist: "Adele", Duration: "4:09"},
	{Title: "Aalilou", Artist: "Amr Diab", Duration: "4:05"},
	{Title: "Nour El عين (Light of My Eye)", Artist: "Amr Diab", Duration: "4:18"},
	{Title: "Rayehn (Leader)", Artist: "Amr Diab", Duration: "4:20"},
	{Title: "I Bet You Look Good on the Dancefloor", Artist: "Arctic Monkeys", Duration: "3:57"},
	{Title: "Do I Wanna Know?", Artist: "Arctic Monkeys", Duration: "4:05"},
	{Title: "Bad Romance", Artist: "Lady Gaga", Duration: "4:54"},
	{Title: "Poker Face", Artist: "Lady Gaga", Duration: "4:03"},
	{Title: "Still D.R.E.", Artist: "Dr Dre", Duration: "4:30"},
	{Title: "California Love", Artist: "Dr Dre", Duration: "4:10"},
	{Title: "Bohemian Rhapsody", Artist: "Queen", Duration: "5:55"},
	{Title: "We Are the Champions", Artist: "Queen", Duration: "3:49"},
	{Title: "Genie in a Bottle", Artist: "Christina Agiulera", Duration: "3:39"},
	{Title: "Dirrty", Artist: "Christina Agiulera", Duration: "4:08"},
	{Title: "Complicated", Artist: "Avril Lavigne", Duration: "4:03"},
	{Title: "Sk8er Boi", Artist: "Avril Lavigne", Duration: "3:33"},
	{Title: "Thank U, Next", Artist: "Ariana Grande", Duration: "3:07"},
	{Title: "7 Rings", Artist: "Ariana Grande", Duration: "2:59"},
	{Title: "99 Problems", Artist: "Jay-Z", Duration: "4:00"},
	{Title: "Empire State of Mind", Artist: "Jay-Z", Duration: "4:22"},
	{Title: "Blinding Lights", Artist: "The Weekend", Duration: "3:40"},
	{Title: "Uptown Funk", Artist: "Mark Ronson ft. Bruno Mars", Duration: "4:30"},
	{Title: "Perfect", Artist: "Ed Sheeran", Duration: "4:02"},
	{Title: "Don't Start Now", Artist: "Dua Lipa", Duration: "3:29"},
	{Title: "Levitating", Artist: "Dua Lipa", Duration: "3:23"},
}

func RandomSong() models.Song {
	song := songs[random.Intn(len(songs))]
	song.YtId = nanoid.Generate()
	return song
}

func RandomSongs(amount int) []*models.Song {
	randSongs := make([]*models.Song, amount)
	for i := 0; i < amount; i++ {
		song := RandomSong()
		randSongs[i] = &song
		random.Seed(time.Now().UnixMicro())
	}
	return randSongs
}

func Songs() []models.Song {
	return songs
}

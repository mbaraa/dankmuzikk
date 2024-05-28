package tests

import (
	"dankmuzikk/models"
	"time"
)

var profiles = []models.Profile{
	{Name: "Shrek"},
	{Name: "Fiona"},
	{Name: "Donkey"},
	{Name: "Lord Farquaad"},
	{Name: "Princess Peach"},
	{Name: "Dragon"},
	{Name: "Puss in Boots"},
	{Name: "Pinocchio"},
	{Name: "Gingy"},
	{Name: "The Three Little Pigs"},
	{Name: "The Wolf"},
	{Name: "Queen Lillian"},
	{Name: "King Harold"},
	{Name: "The Fairy Godmother"},
	{Name: "The Muffin Man"},
	{Name: "Monsieur Hood"},
	{Name: "Robin Hood"},
	{Name: "Little Red Riding Hood"},
	{Name: "The Three Blind Mice"},
	{Name: "Sugar Plum Fairy"},
	{Name: "The Dronkeys"},
	{Name: "Doris"},
	{Name: "Donkey (Dragon)"},
	{Name: "Dragon (Donkey)"},
	{Name: "The Ugly Duckling"},
	{Name: "Prince Charming"},
	{Name: "The Pied Piper"},
	{Name: "Mon Farquaad (Headless)"},
	{Name: "The Gingerbread Man"},
	{Name: "The Cheeseball Vendor"},
	{Name: "The Evil Tree"},
	{Name: "The Knight"},
	{Name: "The Baker's Wife"},
	{Name: "The Baker"},
	{Name: "The Captain of the Guards"},
	{Name: "The Page"},
	{Name: "The Oracle"},
	{Name: "The Matchmaker"},
	{Name: "The Pied Piper's Rats"},
	{Name: "The Merry Men"},
}

func initProfiles() {
	for i := range profiles {
		profiles[i].Account = accounts[i]
		profiles[i].Username = profiles[i].Account.Email
	}
}

func RandomProfile() models.Profile {
	return profiles[random.Intn(len(profiles))]
}

func RandomProfiles(amount int) []models.Profile {
	randProfile := make([]models.Profile, amount)
	for i := 0; i < amount; i++ {
		randProfile[i] = RandomProfile()
		random.Seed(time.Now().UnixMicro())
	}
	return randProfile
}

func Profiles() []models.Profile {
	return profiles
}

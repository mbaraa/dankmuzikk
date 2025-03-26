package layouts

type PageType string

const (
	SongPage     PageType = "music.song"
	AlbumPage    PageType = "music.album"
	PlaylistPage PageType = "music.playlist"
	ProfilePage  PageType = "profile"
)

type AudioProps struct {
	// og:audio
	Url string
	// music:duration
	Duration string
	// music:album
	Album string
	// music:musician
	Musician string
}

type PageProps struct {
	Title       string
	Description string
	Type        PageType
	Url         string
	ImageUrl    string
	Audio       AudioProps
}

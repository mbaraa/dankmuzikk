package actions

type PlaylistPermissions int8

const (
	VisitorPermission PlaylistPermissions = 1 << iota
	JoinerPermission
	OwnerPermission
	NonePermission PlaylistPermissions = 0
)

type Playlist struct {
	PublicId    string              `json:"public_id"`
	Title       string              `json:"title"`
	SongsCount  int                 `json:"songs_count"`
	Songs       []Song              `json:"songs"`
	IsPublic    bool                `json:"is_public"`
	Permissions PlaylistPermissions `json:"permissions"`
}

func (a *Actions) GetAllPlaylists(sessionToken string) ([]Playlist, error) {
	return a.requests.GetPlaylists(sessionToken)
}

func (a *Actions) GetSinglePlaylist(sessionToken, playlistPublicId string) (Playlist, error) {
	return a.requests.GetPlaylist(sessionToken, playlistPublicId)
}

func (a *Actions) CreatePlaylist(sessionToken string, playlist Playlist) (Playlist, error) {
	return a.requests.CreatePlaylist(sessionToken, playlist)
}

func (a *Actions) GetAllPlaylistsForAddPopover(sessionToken string) ([]Playlist, map[string]bool, error) {
	return a.requests.GetAllPlaylistsForAddPopover(sessionToken)
}

func (a *Actions) DeletePlaylist(sessionToken, playlistPublicId string) error {
	return a.requests.DeletePlaylist(sessionToken, playlistPublicId)
}

func (a *Actions) ToggleJoinPlaylist(sessionToken, playlistPublicId string) (joined bool, err error) {
	return a.requests.ToggleJoinPlaylist(sessionToken, playlistPublicId)
}

func (a *Actions) TogglePublicPlaylist(sessionToken, playlistPublicId string) (public bool, err error) {
	return a.requests.TogglePublicPlaylist(sessionToken, playlistPublicId)
}

func (a *Actions) DownloadPlaylist(sessionToken, playlistPublicId string) (string, error) {
	return a.requests.DownloadPlaylist(sessionToken, playlistPublicId)
}

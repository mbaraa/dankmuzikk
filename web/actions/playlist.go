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

func (a *Actions) GetAllPlaylists(ctx ActionContext) ([]Playlist, error) {
	return a.requests.GetPlaylists(ctx.SessionToken)
}

type GetSinglePlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) GetSinglePlaylist(params GetSinglePlaylistParams) (Playlist, error) {
	return a.requests.GetPlaylist(params.SessionToken, params.PlaylistPublicId)
}

type CreatePlaylistParams struct {
	ActionContext
	Playlist Playlist
}

func (a *Actions) CreatePlaylist(params CreatePlaylistParams) (Playlist, error) {
	return a.requests.CreatePlaylist(params.SessionToken, params.Playlist)
}

func (a *Actions) GetAllPlaylistsForAddPopover(ctx ActionContext) ([]Playlist, map[string]bool, error) {
	return a.requests.GetAllPlaylistsForAddPopover(ctx.SessionToken)
}

type DeletePlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) DeletePlaylist(params DeletePlaylistParams) error {
	return a.requests.DeletePlaylist(params.SessionToken, params.PlaylistPublicId)
}

type ToggleJoinPlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) ToggleJoinPlaylist(params ToggleJoinPlaylistParams) (joined bool, err error) {
	return a.requests.ToggleJoinPlaylist(params.SessionToken, params.PlaylistPublicId)
}

type TogglePublicPlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) TogglePublicPlaylist(params TogglePublicPlaylistParams) (public bool, err error) {
	return a.requests.TogglePublicPlaylist(params.SessionToken, params.PlaylistPublicId)
}

type DownloadPlaylistParams struct {
	ActionContext
	PlaylistPublicId string
}

func (a *Actions) DownloadPlaylist(params DownloadPlaylistParams) (string, error) {
	return a.requests.DownloadPlaylist(params.SessionToken, params.PlaylistPublicId)
}

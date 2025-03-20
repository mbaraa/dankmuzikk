package entities

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

package playlists

import "errors"

var (
	ErrOwnerCantLeavePlaylist      = errors.New("playlists: owner can't leave playlists")
	ErrNonOwnerCantDeletePlaylists = errors.New("playlists: non owners can only leave playlists")
)

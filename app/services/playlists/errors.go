package playlists

import "errors"

var (
	ErrOwnerCantLeavePlaylist      = errors.New("playlists: owner can't leave playlists")
	ErrNonOwnerCantDeletePlaylists = errors.New("playlists: non owners can only leave playlists")
	ErrUnauthorizedToSeePlaylist   = errors.New("playlists: unauthorized to see playlist")
	ErrEmptyPlaylist               = errors.New("playlists: empty playlists")
	ErrUserHasAlreadyVoted         = errors.New("playlist: user can't vote more than once")
)

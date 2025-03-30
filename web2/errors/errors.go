package errors

import "errors"

var (
	ErrInvalidToken                           = errors.New("invalid-token")
	ErrExpiredToken                           = errors.New("expired-token")
	ErrAccountNotFound                        = errors.New("account-not-found")
	ErrProfileNotFound                        = errors.New("profile-not-found")
	ErrSongNotFound                           = errors.New("song-not-found")
	ErrPlaylistNotFound                       = errors.New("playlist-not-found")
	ErrAccountExists                          = errors.New("account-exists")
	ErrProfileExists                          = errors.New("profile-exists")
	ErrSongExists                             = errors.New("song-exists")
	ErrPlaylistExists                         = errors.New("playlist-exists")
	ErrDifferentLoginMethodUsed               = errors.New("different-login-method-used")
	ErrVerificationCodeExpired                = errors.New("verification-code-expired")
	ErrInvalidVerificationCode                = errors.New("invalid-verification-code")
	ErrNonOwnerCantDeletePlaylists            = errors.New("non-owner-cant-delete-playlists")
	ErrNonOwnerCantChangePlaylistVisibility   = errors.New("non-owner-cant-change-playlist-visibility")
	ErrUnauthorizedToSeePlaylist              = errors.New("unauthorized-to-see-playlist")
	ErrUserHasAlreadyVoted                    = errors.New("user-has-already-voted")
	ErrNotEnoughPermissionToAddSongToPlaylist = errors.New("not-enough-permission-to-add-song-to-playlist")
	ErrInvalidSessionToken                    = errors.New("invalid-session-token")

	ErrSomethingWentWrong = errors.New("something went wrong")
)

package requests

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

	ErrSomethingWentWrong = errors.New("something went wrong")
)

func mapError(errorId string) error {
	switch errorId {
	case "invalid-token":
		return ErrInvalidToken
	case "expired-token":
		return ErrExpiredToken
	case "account-not-found":
		return ErrAccountNotFound
	case "profile-not-found":
		return ErrProfileNotFound
	case "song-not-found":
		return ErrSongNotFound
	case "playlist-not-found":
		return ErrPlaylistNotFound
	case "account-exists":
		return ErrAccountExists
	case "profile-exists":
		return ErrProfileExists
	case "song-exists":
		return ErrSongExists
	case "playlist-exists":
		return ErrPlaylistExists
	case "different-login-method-used":
		return ErrDifferentLoginMethodUsed
	case "verification-code-expired":
		return ErrVerificationCodeExpired
	case "invalid-verification-code":
		return ErrInvalidVerificationCode
	case "non-owner-cant-delete-playlists":
		return ErrNonOwnerCantDeletePlaylists
	case "non-owner-cant-change-playlist-visibility":
		return ErrNonOwnerCantChangePlaylistVisibility
	case "unauthorized-to-see-playlist":
		return ErrUnauthorizedToSeePlaylist
	case "user-has-already-voted":
		return ErrUserHasAlreadyVoted
	case "not-enough-permission-to-add-song-to-playlist":
		return ErrNotEnoughPermissionToAddSongToPlaylist
	default:
		return ErrSomethingWentWrong
	}
}

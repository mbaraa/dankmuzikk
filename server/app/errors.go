package app

import (
	"fmt"
	"net/http"
	"strings"
)

// DankError is implemented for every error around here :)
type DankError interface {
	error
	// ClientStatusCode the HTTP status for clients.
	ClientStatusCode() int
	// ExtraData any data that will be helpful for clients for better UX context.
	ExtraData() map[string]any
	// ExposeToClients reports whether to expose this error to clients or not.
	ExposeToClients() bool
}

type ErrNotFound struct {
	ResourceName string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s-not-found", strings.ToLower(e.ResourceName))
}

func (e ErrNotFound) ClientStatusCode() int {
	return http.StatusNotFound
}

func (e ErrNotFound) ExtraData() map[string]any {
	return nil
}

func (e ErrNotFound) ExposeToClients() bool {
	return true
}

type ErrExists struct {
	ResourceName string
}

func (e ErrExists) Error() string {
	return fmt.Sprintf("%s-exists", strings.ToLower(e.ResourceName))
}

func (e ErrExists) ClientStatusCode() int {
	return http.StatusConflict
}

func (e ErrExists) ExtraData() map[string]any {
	return nil
}

func (e ErrExists) ExposeToClients() bool {
	return true
}

type ErrDifferentLoginMethod struct{}

func (e ErrDifferentLoginMethod) Error() string {
	return "different-login-method-used"
}

func (e ErrDifferentLoginMethod) ClientStatusCode() int {
	return http.StatusConflict
}

func (e ErrDifferentLoginMethod) ExtraData() map[string]any {
	return nil
}

func (e ErrDifferentLoginMethod) ExposeToClients() bool {
	return true
}

type ErrExpiredVerificationCode struct{}

func (e ErrExpiredVerificationCode) Error() string {
	return "verification-code-expired"
}

func (e ErrExpiredVerificationCode) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrExpiredVerificationCode) ExtraData() map[string]any {
	return nil
}

func (e ErrExpiredVerificationCode) ExposeToClients() bool {
	return true
}

type ErrInvalidSessionToken struct{}

func (e ErrInvalidSessionToken) Error() string {
	return "invalid-session-code"
}

func (e ErrInvalidSessionToken) ClientStatusCode() int {
	return http.StatusBadRequest
}

func (e ErrInvalidSessionToken) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidSessionToken) ExposeToClients() bool {
	return true
}

type ErrInvalidVerificationToken struct{}

func (e ErrInvalidVerificationToken) Error() string {
	return "invalid-verification-code"
}

func (e ErrInvalidVerificationToken) ClientStatusCode() int {
	return http.StatusBadRequest
}

func (e ErrInvalidVerificationToken) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidVerificationToken) ExposeToClients() bool {
	return true
}

type ErrNonOwnerCantDeletePlaylists struct{}

func (e ErrNonOwnerCantDeletePlaylists) Error() string {
	return "non-owner-cant-delete-playlists"
}

func (e ErrNonOwnerCantDeletePlaylists) ClientStatusCode() int {
	return http.StatusForbidden
}

func (e ErrNonOwnerCantDeletePlaylists) ExtraData() map[string]any {
	return nil
}

func (e ErrNonOwnerCantDeletePlaylists) ExposeToClients() bool {
	return true
}

type ErrNonOwnerCantChangePlaylistVisibility struct{}

func (e ErrNonOwnerCantChangePlaylistVisibility) Error() string {
	return "non-owner-cant-change-playlist-visibility"
}

func (e ErrNonOwnerCantChangePlaylistVisibility) ClientStatusCode() int {
	return http.StatusForbidden
}

func (e ErrNonOwnerCantChangePlaylistVisibility) ExtraData() map[string]any {
	return nil
}

func (e ErrNonOwnerCantChangePlaylistVisibility) ExposeToClients() bool {
	return true
}

type ErrUnauthorizedToSeePlaylist struct{}

func (e ErrUnauthorizedToSeePlaylist) Error() string {
	return "unauthorized-to-see-playlist"
}

func (e ErrUnauthorizedToSeePlaylist) ClientStatusCode() int {
	return http.StatusForbidden
}

func (e ErrUnauthorizedToSeePlaylist) ExtraData() map[string]any {
	return nil
}

func (e ErrUnauthorizedToSeePlaylist) ExposeToClients() bool {
	return true
}

type ErrUserHasAlreadyVoted struct{}

func (e ErrUserHasAlreadyVoted) Error() string {
	return "user-has-already-voted"
}

func (e ErrUserHasAlreadyVoted) ClientStatusCode() int {
	return http.StatusConflict
}

func (e ErrUserHasAlreadyVoted) ExtraData() map[string]any {
	return nil
}

func (e ErrUserHasAlreadyVoted) ExposeToClients() bool {
	return true
}

type ErrNotEnoughPermissionToAddSongToPlaylist struct{}

func (e ErrNotEnoughPermissionToAddSongToPlaylist) Error() string {
	return "not-enough-permission-to-add-song-to-playlist"
}

func (e ErrNotEnoughPermissionToAddSongToPlaylist) ClientStatusCode() int {
	return http.StatusForbidden
}

func (e ErrNotEnoughPermissionToAddSongToPlaylist) ExtraData() map[string]any {
	return nil
}

func (e ErrNotEnoughPermissionToAddSongToPlaylist) ExposeToClients() bool {
	return true
}

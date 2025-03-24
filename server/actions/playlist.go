package actions

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CreatePlaylistParams struct {
	Title     string `json:"title"`
	ProfileId uint   `json:"-"`
}

type Playlist struct {
	PublicId    string                     `json:"public_id"`
	Title       string                     `json:"title"`
	SongsCount  int                        `json:"songs_count"`
	Songs       []Song                     `json:"songs"`
	IsPublic    bool                       `json:"is_public"`
	Permissions models.PlaylistPermissions `json:"permissions"`
}

type CreatePlaylistPayload struct {
	NewPlaylist Playlist `json:"new_playlist"`
}

func (a *Actions) CreatePlaylist(params CreatePlaylistParams) (CreatePlaylistPayload, error) {
	playlist, err := a.app.CreatePlaylist(app.CreatePlaylistArgs{
		Title:     params.Title,
		ProfileId: params.ProfileId,
	})
	if err != nil {
		return CreatePlaylistPayload{}, err
	}

	return CreatePlaylistPayload{
		NewPlaylist: Playlist{
			PublicId:   playlist.PublicId,
			Title:      playlist.Title,
			SongsCount: playlist.SongsCount,
			IsPublic:   playlist.IsPublic,
		},
	}, nil
}

func (a *Actions) TogglePublicPlaylist(playlistPubId string, ownerId uint) (madePublic bool, err error) {
	return a.app.TogglePublicPlaylist(playlistPubId, ownerId)
}

func (a *Actions) ToggleProfileInPlaylist(playlistPubId string, profileId uint) (joined bool, err error) {
	return a.app.ToggleProfileInPlaylist(playlistPubId, profileId)
}

func (a *Actions) GetPlaylistByPublicId(playlistPubId string, profileId uint) (Playlist, error) {
	playlist, permissions, err := a.app.GetPlaylistByPublicId(playlistPubId, profileId)
	if err != nil {
		return Playlist{}, err
	}

	var songs []Song
	for _, song := range playlist.Songs {
		songs = append(songs, Song{
			YtId:         song.YtId,
			Title:        song.Title,
			Artist:       song.Artist,
			ThumbnailUrl: song.ThumbnailUrl,
			Duration:     song.Duration,
			PlayTimes:    song.PlayTimes,
			Votes:        song.Votes,
			AddedAt:      song.AddedAt,
		})
	}

	return Playlist{
		PublicId:    playlist.PublicId,
		Title:       playlist.Title,
		SongsCount:  playlist.SongsCount,
		Songs:       songs,
		IsPublic:    playlist.IsPublic,
		Permissions: permissions,
	}, nil
}

func (a *Actions) DeletePlaylist(playlistPubId string, profileId uint) error {
	return a.app.DeletePlaylist(playlistPubId, profileId)
}

func (a *Actions) DownloadPlaylist(playlistPubId string, profileId uint) (string, error) {
	playlist, _, err := a.app.GetPlaylistByPublicId(playlistPubId, profileId)
	if err != nil {
		return "", err
	}

	fileNames := make([]string, 0, playlist.SongsCount)
	for i, song := range playlist.Songs {
		ogFile, err := os.Open(fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().BlobsDir, song.YtId))
		if err != nil {
			return "", err
		}
		newShit, err := os.OpenFile(
			filepath.Clean(
				fmt.Sprintf("%s/muzikkx/%d-%s.mp3", config.Env().BlobsDir, i+1,
					strings.ReplaceAll(song.Title, "/", "|"),
				),
			),
			os.O_WRONLY|os.O_CREATE, 0644,
		)
		if err != nil {
			_ = ogFile.Close()
			return "", err
		}
		_, err = io.Copy(newShit, ogFile)
		if err != nil {
			_ = ogFile.Close()
			return "", err
		}
		fileNames = append(fileNames, newShit.Name())
		_ = newShit.Close()
		_ = ogFile.Close()
	}

	archive, err := a.archiver.CreateArchive(playlist.Title)
	if err != nil {
		return "", err
	}

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			return "", err
		}
		err = archive.AddFile(file)
		if err != nil {
			return "", err
		}
		_ = file.Close()
		_ = os.Remove(file.Name())
	}

	playlistZip, err := archive.Deflate()
	if err != nil {
		return "", err
	}

	outFile, err := os.Create(fmt.Sprintf("%s/playlists/%s.zip", config.Env().BlobsDir, playlist.PublicId))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(outFile, playlistZip)
	if err != nil {
		return "", err
	}

	_ = outFile.Close()

	return fmt.Sprintf("%s/playlists/%s.zip", config.Env().CdnAddress, playlist.PublicId), nil
}

func (a *Actions) DeletePlaylistArchive(event events.PlaylistDownloaded) error {
	if event.DeleteAt.Before(time.Now().UTC()) {
		return a.eventhub.Publish(event)
	}

	err := os.Remove(fmt.Sprintf("%s/playlists/%s.zip", config.Env().BlobsDir, event.PlaylistId))
	if err != nil {
		return err
	}

	return nil
}

func (a *Actions) GetAllPlaylistsMappedWithSongs(ownerId uint) ([]Playlist, map[string]bool, error) {
	playlists, mapping, err := a.app.GetAllPlaylistsMappedWithSongs(ownerId)
	if err != nil {
		return nil, nil, err
	}

	outPlaylists := make([]Playlist, 0, len(playlists))
	for _, playlist := range playlists {
		outPlaylists = append(outPlaylists, Playlist{
			PublicId:   playlist.PublicId,
			Title:      playlist.Title,
			SongsCount: playlist.SongsCount,
			IsPublic:   playlist.IsPublic,
		})
	}

	return outPlaylists, mapping, nil
}

func (a *Actions) GetPlaylistsForProfile(ownerId uint) ([]Playlist, error) {
	playlists, err := a.app.GetPlaylistsForProfile(ownerId)
	if err != nil {
		return nil, err
	}

	outPlaylists := make([]Playlist, 0, playlists.Size)
	for playlist := range playlists.Seq() {
		outPlaylists = append(outPlaylists, Playlist{
			PublicId:   playlist.PublicId,
			Title:      playlist.Title,
			SongsCount: playlist.SongsCount,
			IsPublic:   playlist.IsPublic,
		})
	}

	return outPlaylists, nil
}

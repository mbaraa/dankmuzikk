package actions

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CreatePlaylistParams struct {
	Title     string `json:"title"`
	ProfileId uint   `json:"-"`
}

type Playlist struct {
	PublicId   string `json:"public_id"`
	Title      string `json:"title"`
	SongsCount int    `json:"songs_count"`
	Songs      []Song `json:"songs"`
	IsPublic   bool   `json:"is_public"`
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
	playlist, err := a.app.GetPlaylistByPublicId(playlistPubId, profileId)
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
		PublicId:   playlist.PublicId,
		Title:      playlist.Title,
		SongsCount: playlist.SongsCount,
		Songs:      songs,
		IsPublic:   playlist.IsPublic,
	}, nil
}

func (a *Actions) DeletePlaylist(playlistPubId string, profileId uint) error {
	return a.app.DeletePlaylist(playlistPubId, profileId)
}

func (a *Actions) DownloadPlaylist(playlistPubId string, profileId uint) (io.Reader, error) {
	playlist, err := a.app.GetPlaylistByPublicId(playlistPubId, profileId)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, playlist.SongsCount)
	for i, song := range playlist.Songs {
		ogFile, err := os.Open(fmt.Sprintf("%s/%s.mp3", config.Env().YouTube.MuzikkDir, song.YtId))
		if err != nil {
			return nil, err
		}
		newShit, err := os.OpenFile(
			filepath.Clean(
				fmt.Sprintf("%s/%d-%s.mp3", config.Env().YouTube.MuzikkDir, i+1,
					strings.ReplaceAll(song.Title, "/", "|"),
				),
			),
			os.O_WRONLY|os.O_CREATE, 0644,
		)
		if err != nil {
			_ = ogFile.Close()
			return nil, err
		}
		_, err = io.Copy(newShit, ogFile)
		if err != nil {
			_ = ogFile.Close()
			return nil, err
		}
		fileNames[i] = newShit.Name()
		_ = newShit.Close()
		_ = ogFile.Close()
	}

	archive, err := a.archiver.CreateArchive(playlist.Title)
	if err != nil {
		return nil, err
	}

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		err = archive.AddFile(file)
		if err != nil {
			return nil, err
		}
		_ = file.Close()
		_ = os.Remove(file.Name())
	}

	defer func() {
	}()

	return archive.Deflate()
}

func (a *Actions) GetAllPlaylistsMappedWithSongs(ownerId uint) ([]models.Playlist, map[string]bool, error) {
	return a.app.GetAllPlaylistsMappedWithSongs(ownerId)
}

func (a *Actions) GetPlaylistsForProfile(ownerId uint) (models.List[Playlist], error) {
	playlists, err := a.app.GetPlaylistsForProfile(ownerId)
	if err != nil {
		return models.List[Playlist]{}, err
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

	return models.NewList(outPlaylists, ""), nil
}

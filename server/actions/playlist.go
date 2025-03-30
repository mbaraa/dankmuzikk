package actions

import (
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"dankmuzikk/evy/events"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Playlist struct {
	PublicId    string                     `json:"public_id"`
	Title       string                     `json:"title"`
	SongsCount  int                        `json:"songs_count"`
	Songs       []Song                     `json:"songs"`
	IsPublic    bool                       `json:"is_public"`
	Permissions models.PlaylistPermissions `json:"permissions"`
}

type CreatePlaylistParams struct {
	ActionContext `json:"-"`
	Title         string `json:"title"`
}

type CreatePlaylistPayload struct {
	NewPlaylist Playlist `json:"new_playlist"`
}

func (a *Actions) CreatePlaylist(params CreatePlaylistParams) (CreatePlaylistPayload, error) {
	playlist, err := a.app.CreatePlaylist(app.CreatePlaylistArgs{
		Title:     params.Title,
		AccountId: params.Account.Id,
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

type TogglePublicPlaylistParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string
}

type TogglePublicPlaylistPayload struct {
	ActionContext `json:"-"`
	Public        bool `json:"public"`
}

func (a *Actions) TogglePublicPlaylist(params TogglePublicPlaylistParams) (TogglePublicPlaylistPayload, error) {
	madePublic, err := a.app.TogglePublicPlaylist(params.PlaylistPublicId, params.Account.Id)
	if err != nil {
		return TogglePublicPlaylistPayload{}, err
	}

	return TogglePublicPlaylistPayload{
		Public: madePublic,
	}, nil

}

type ToggleJoinPlaylistParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string
}

type ToggleJoinPlaylistPayload struct {
	Joined bool `json:"joined"`
}

func (a *Actions) ToggleJoinPlaylist(params ToggleJoinPlaylistParams) (ToggleJoinPlaylistPayload, error) {
	joined, err := a.app.ToggleAccountInPlaylist(params.PlaylistPublicId, params.Account.Id)
	if err != nil {
		return ToggleJoinPlaylistPayload{}, err
	}

	return ToggleJoinPlaylistPayload{
		Joined: joined,
	}, nil
}

type GetPlaylistByPublicIdParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string
}

func (a *Actions) GetPlaylistByPublicId(params GetPlaylistByPublicIdParams) (Playlist, error) {
	playlist, permissions, err := a.app.GetPlaylistByPublicId(params.PlaylistPublicId, params.Account.Id)
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

type DeletePlaylistParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string
}

func (a *Actions) DeletePlaylist(params DeletePlaylistParams) error {
	return a.app.DeletePlaylist(params.PlaylistPublicId, params.Account.Id)
}

type DownloadPlaylistParams struct {
	ActionContext    `json:"-"`
	PlaylistPublicId string
}

type DownloadPlaylistPayload struct {
	PlaylistDownloadUrl string `json:"playlist_download_url"`
}

func (a *Actions) DownloadPlaylist(params DownloadPlaylistParams) (DownloadPlaylistPayload, error) {
	playlist, _, err := a.app.GetPlaylistByPublicId(params.PlaylistPublicId, params.Account.Id)
	if err != nil {
		return DownloadPlaylistPayload{}, err
	}

	fileNames := make([]string, 0, playlist.SongsCount)
	var errs []error
	for i, song := range playlist.Songs {
		oldPath := fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().BlobsDir, song.YtId)
		newPath := fmt.Sprintf("%s/muzikkx/%d-%s.mp3", config.Env().BlobsDir, i+1,
			strings.ReplaceAll(song.Title, "/", "|"),
		)
		err = a.blobstorage.CopyFile(oldPath, newPath)
		if err != nil {
			// TODO: add blobstorage.Exists or something, just to make sure the error is coming from one place.
			errs = append(errs, err)
			// to download the song, so a second download retry would be enough without playing all the songs :)
			_ = a.eventhub.Publish(events.SongPlayed{
				SongYtId: song.YtId,
			})
			continue
		}

		fileNames = append(fileNames, newPath)
	}
	if len(errs) != 0 {
		return DownloadPlaylistPayload{}, errors.Join(errs...)
	}

	archive, err := a.archiver.CreateArchive(playlist.Title)
	if err != nil {
		return DownloadPlaylistPayload{}, err
	}

	for _, fileName := range fileNames {
		file, err := a.blobstorage.GetFile(fileName)
		if err != nil {
			return DownloadPlaylistPayload{}, err
		}
		err = archive.AddFile(file)
		if err != nil {
			return DownloadPlaylistPayload{}, err
		}

		_ = a.blobstorage.DeleteFile(file.Name())
		_ = file.Close()
	}

	playlistZip, err := archive.Deflate()
	if err != nil {
		return DownloadPlaylistPayload{}, err
	}

	playlistsArchivePath := fmt.Sprintf("%s/playlists/%s.zip", config.Env().BlobsDir, playlist.PublicId)
	err = a.blobstorage.CreateFile(playlistsArchivePath)
	if err != nil {
		return DownloadPlaylistPayload{}, err
	}

	err = a.blobstorage.WriteToFile(playlistsArchivePath, playlistZip)
	if err != nil {
		return DownloadPlaylistPayload{}, err
	}

	return DownloadPlaylistPayload{
		PlaylistDownloadUrl: fmt.Sprintf("%s/playlists/%s.zip", config.Env().CdnAddress, playlist.PublicId),
	}, nil
}

type GetAllPlaylistsMappedWithSongsParams struct {
	ActionContext `json:"-"`
}

type GetAllPlaylistsMappedWithSongsPayload struct {
	Playlists []Playlist `json:"playlists"`
	// TODO: maybe just send the playlists mapped :)
	SongsInPlaylists map[string]bool `json:"songs_in_playlists"`
}

func (a *Actions) GetAllPlaylistsMappedWithSongs(params GetAllPlaylistsMappedWithSongsParams) (GetAllPlaylistsMappedWithSongsPayload, error) {
	playlists, mapping, err := a.app.GetAllPlaylistsMappedWithSongs(params.Account.Id)
	if err != nil {
		return GetAllPlaylistsMappedWithSongsPayload{}, err
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

	return GetAllPlaylistsMappedWithSongsPayload{
		Playlists:        outPlaylists,
		SongsInPlaylists: mapping,
	}, nil
}

type GetPlaylistsForAccountParams struct {
	ActionContext `json:"-"`
}

// TODO: use this
type GetPlaylistsForAccountPayload struct {
	Data []Playlist `json:"data"`
}

func (a *Actions) GetPlaylistsForAccount(params GetPlaylistsForAccountParams) ([]Playlist, error) {
	playlists, err := a.app.GetPlaylistsForAccount(params.Account.Id)
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

func (a *Actions) DeletePlaylistArchive(event events.PlaylistDownloaded) error {
	if event.DeleteAt.Before(time.Now().UTC()) {
		return a.eventhub.Publish(event)
	}

	err := a.blobstorage.DeleteFile(fmt.Sprintf("%s/playlists/%s.zip", config.Env().BlobsDir, event.PlaylistId))
	if err != nil {
		return err
	}

	return nil
}

func (a *Actions) IncrementSongsCountForPlaylist(playlistPublicId string, accountId uint) error {
	return a.app.IncrementSongsCountForPlaylist(playlistPublicId, accountId)
}

func (a *Actions) DecrementSongsCountForPlaylist(playlistPublicId string, accountId uint) error {
	return a.app.DecrementSongsCountForPlaylist(playlistPublicId, accountId)
}

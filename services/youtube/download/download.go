package download

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
)

// Service represents the YouTube downloader service.
type Service struct {
	repo db.CRUDRepo[models.Song]
}

// New accepts a models.Song creator repo to store songs' meta-data,
// and returns a new instance of the YouTube downloader service.
func New(repo db.CRUDRepo[models.Song]) *Service {
	return &Service{repo}
}

// DownloadYoutubeSong downloads a YouTube music file into the path specified by the environment variable
// YOUTUBE_MUSIC_DOWNLOAD_PATH, where the file name will be <video_id.mp3> to be served under /music/{id}
// and returns an occurring error
//
// Used when playing a new song (usually from search).
// TODO: optimize select query, maybe?
func (d *Service) DownloadYoutubeSong(songYtId string) error {
	song, err := d.repo.GetByConds("yt_id = ?", songYtId)
	if err == nil && len(song) != 0 && song[0].FullyDownloaded {
		log.Infof("The song with id %s is already downloaded\n", songYtId)
		return nil
	}

	err = d.DownloadYoutubeSongQueue(songYtId)
	if err != nil {
		return err
	}
	resp, err := http.Get(fmt.Sprintf("%s/download/%s", config.Env().YouTube.DownloaderUrl, songYtId))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong when downloading a song; id: " + songYtId)
	}

	return nil
}

// DownloadYoutubeSongQueue same as DownloadYoutubeSong but it downloads the song in the background,
// and only downloads the song's file (since the meta was already downloaded before).
//
// Used when adding a song to a playlist.
func (d *Service) DownloadYoutubeSongQueue(songYtId string) error {
	resp, err := http.Get(fmt.Sprintf("%s/download/queue/%s", config.Env().YouTube.DownloaderUrl, songYtId))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong when downloading a song; id: " + songYtId)
	}

	return nil
}

// DownloadYoutubeSongsMetadata same as DownloadYoutubeSong but it downloads the song in the background
//
// Used when searching for a query.
// TODO: move this logic out of here, since it doesn't fit in here :)
func (d *Service) DownloadYoutubeSongsMetadata(songs []entities.Song) error {
	newSongs := make([]*models.Song, 0)
	for _, song := range songs {
		path := fmt.Sprintf("%s/%s.mp3", config.Env().YouTube.MusicDir, song.YtId)
		if _, err := os.Stat(path); err == nil {
			log.Infof("The song with id %s is already downloaded\n", song.YtId)
			continue
		}

		newSongs = append(newSongs, &models.Song{
			YtId:         song.YtId,
			Title:        song.Title,
			Artist:       song.Artist,
			ThumbnailUrl: song.ThumbnailUrl,
			Duration:     song.Duration,
		})
	}

	// adding the songs one at a time, so that if a song exists, it won't ruin the batch!
	for _, newSong := range newSongs {
		err := d.repo.Add(newSong)
		if errors.Is(err, db.ErrRecordExists) {
			log.Warningln(err)
		} else if err != nil {
			return err
		}
	}

	return nil
}

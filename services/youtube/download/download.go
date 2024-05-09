package download

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Service represents the YouTube downloader service.
type Service struct {
	repo db.CreatorRepo[models.Song]
}

// New accepts a models.Song creator repo to store songs' meta-data,
// and returns a new instance of the YouTube downloader service.
func New(repo db.CreatorRepo[models.Song]) *Service {
	return &Service{repo}
}

// DownloadYoutubeSong downloads a YouTube music file into the path specified by the environment variable
// YOUTUBE_MUSIC_DOWNLOAD_PATH, where the file name will be <video_id.mp3> to be served under /music/{id}
// and returns an occurring error
func (d *Service) DownloadYoutubeSong(req entities.SongDownloadRequest) error {
	path := fmt.Sprintf("%s/%s.mp3", config.Env().YouTube.MusicDir, req.Id)
	if _, err := os.Stat(path); err == nil {
		log.Infof("The song with id %s is already downloaded\n", req.Id)
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("%s/download/%s", config.Env().YouTube.DownloaderUrl, req.Id))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong when downloading a song; id: " + req.Id)
	}

	err = d.repo.Add(&models.Song{
		YtId:         req.Id,
		Title:        req.Title,
		Artist:       req.Artist,
		ThumbnailUrl: req.ThumbnailUrl,
		Duration:     req.Duration,
	})
	if err != nil {
		return err
	}

	return nil
}

package download

import (
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type DownloadService struct {
	repo db.CreatorRepo[models.Song]
}

func NewDownloadService(repo db.CRUDRepo[models.Song]) *DownloadService {
	return &DownloadService{repo}
}

// DownloadYoutubeVideo downloads a youtube music file into the path specified by the environment variable
// YOUTUBE_MUSIC_DOWNLOAD_PATH, where the file name will be <video_id.mp3> to be served under /music/{id}
// and retuens an occurring error
func (d *DownloadService) DownloadYoutubeSong(req entities.SongDownloadRequest) error {
	log.Infoln(req)
	path := fmt.Sprintf("%s/%s.mp3", config.Env().YouTube.MusicDir, req.Id)
	if _, err := os.Stat(path); err == nil {
		log.Infof("The song with id %s is already downloaded\n", req.Id)
		return nil
	}
	// TODO: write a downloader instead of using this sub process command thingy.
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--audio-quality", "best", "-o", path, "https://www.youtube.com/watch?v="+req.Id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("Download failed:" + err.Error())
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

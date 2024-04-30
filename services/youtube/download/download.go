package download

import (
	"dankmuzikk/config"
	"dankmuzikk/log"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// DownloadYoutubeVideo downloads a youtube music file into the path specified by the environment variable
// YOUTUBE_MUSIC_DOWNLOAD_PATH, where the file name will be <video_id.mp3> to be served under /music/{id}
// and retuens an occurring error
func DownloadYoutubeVideo(id string) error {
	path := fmt.Sprintf("%s/%s.mp3", config.Vals().YouTube.MusicDir, id)
	if _, err := os.Stat(path); err == nil {
		log.Infof("The song with id %s is already downloaded\n", id)
		return nil
	}
	// TODO: write a downloader instead of using this sub process command thingy.
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--audio-quality", "best", "-o", path, "https://www.youtube.com/watch?v="+id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("Download failed:" + err.Error())
	}
	return nil
}

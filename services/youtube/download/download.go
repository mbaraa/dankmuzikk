package download

import (
	"dankmuzikk/log"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func init() {
	if os.Getenv("YOUTUBE_MUSIC_DOWNLOAD_PATH") == "" {
		log.Fatalln(log.ErrorLevel, "[YOUTUBE DOWNLOAD SERVICE] Missing YouTube Music Download Path")
	}
}

// DownloadYoutubeVideo downloads a youtube music file into the path specified by the environment variable
// YOUTUBE_MUSIC_DOWNLOAD_PATH, where the file name will be <video_id.mp3> to be served under /music/{id}
// and retuens an occurring error
func DownloadYoutubeVideo(id string) error {
	path := os.Getenv("YOUTUBE_MUSIC_DOWNLOAD_PATH")
	// TODO: check if the file is already downloaded!
	// TODO: write a downloader instead of using this sub process command thingy.
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "--audio-quality", "best", "-o", fmt.Sprintf("%s/%s.%%(ext)s", path, id), "https://www.youtube.com/watch?v="+id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("Download failed:" + err.Error())
	}
	return nil
}

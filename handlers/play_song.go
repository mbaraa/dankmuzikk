package handlers

import (
	"dankmuzikk/log"
	"dankmuzikk/services/youtube/download"
	"net/http"
	"os"
)

func HandleServeSongs(hand *http.ServeMux) {
	hand.Handle("/music/", http.StripPrefix("/music", http.FileServer(http.Dir(os.Getenv("YOUTUBE_MUSIC_DOWNLOAD_PATH")))))
}

func HandleDownloadSong(hand *http.ServeMux) {
	hand.HandleFunc("/api/song/download/{youtube_video_id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("youtube_video_id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := download.DownloadYoutubeVideo(id)
		if err != nil {
			log.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

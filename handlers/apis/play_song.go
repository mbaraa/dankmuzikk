package apis

import (
	"dankmuzikk/log"
	"dankmuzikk/services/youtube/download"
	"net/http"
)

func HandleDownloadSong(w http.ResponseWriter, r *http.Request) {
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
}

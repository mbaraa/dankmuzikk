package entities

type SongDownloadRequest struct {
	Id           string `json:"id"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	Duration     string `json:"duration"`
}

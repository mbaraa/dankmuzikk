package entities

type Song struct {
	YtId         string `json:"yt_id"`
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Duration     string `json:"duration"`
	PlayTimes    int    `json:"play_times"`
}

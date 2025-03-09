package actions

import (
	"fmt"
)

func (a *Actions) GetHistoryItems(profileId, page uint) ([]Song, error) {
	songs, err := a.app.GetHistory(profileId, page)
	if err != nil {
		return nil, err
	}

	songsFr := make([]Song, 0, songs.Size)
	for i, song := range songs.Seq2() {
		playTimes := 1
		for ; i < songs.Size-1 && songs.Items[i+1].YtId == songs.Items[i].YtId; i++ {
			playTimes++
		}
		songsFr = append(songsFr, Song{
			YtId:         song.YtId,
			Title:        song.Title,
			Artist:       song.Artist,
			ThumbnailUrl: song.ThumbnailUrl,
			Duration:     song.Duration,
			// whatever that is :)
			AddedAt: fmt.Sprintf("Played %s - %s", times(playTimes), songs.Items[i].AddedAt),
		})
	}

	return songsFr, nil
}

func times(times int) string {
	switch {
	case times == 1:
		return "once"
	case times > 1:
		return fmt.Sprintf("%d times", times)
	default:
		return ""
	}
}

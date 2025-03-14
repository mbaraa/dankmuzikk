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
	for i := 0; i < songs.Size; i++ {
		playTimes := 1
		for ; i < songs.Size-1 && songs.Items[i+1].YtId == songs.Items[i].YtId; i++ {
			playTimes++
		}
		songsFr = append(songsFr, Song{
			YtId:         songs.Items[i].YtId,
			Title:        songs.Items[i].Title,
			Artist:       songs.Items[i].Artist,
			ThumbnailUrl: songs.Items[i].ThumbnailUrl,
			Duration:     songs.Items[i].Duration,
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

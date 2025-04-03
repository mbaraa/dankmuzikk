package actions

import (
	"fmt"
)

type GetHistoryItemsParams struct {
	ActionContext `json:"-"`
	PageIndex     uint `json:"page_index"`
}

// TODO: use this
type GetHistoryItemsPayload struct {
	Data []Song `json:"data"`
}

func (a *Actions) GetHistoryItems(params GetHistoryItemsParams) ([]Song, error) {
	songs, err := a.app.GetHistory(params.Account.Id, params.PageIndex)
	if err != nil {
		return nil, err
	}

	songsFr := make([]Song, 0, songs.Size)
	for i := 0; i < songs.Size; i++ {
		playTimes := 1
		for ; i < songs.Size-1 && songs.Items[i+1].PublicId == songs.Items[i].PublicId; i++ {
			playTimes++
		}
		songsFr = append(songsFr, Song{
			PublicId:     songs.Items[i].PublicId,
			Title:        songs.Items[i].Title,
			Artist:       songs.Items[i].Artist,
			ThumbnailUrl: songs.Items[i].ThumbnailUrl,
			Duration:     songs.Items[i].RealDuration,
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

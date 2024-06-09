package history

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"fmt"
	"time"
)

type Service struct {
	repo     db.UnsafeCRUDRepo[models.History]
	songRepo db.GetterRepo[models.Song]
}

func New(repo db.UnsafeCRUDRepo[models.History], songRepo db.GetterRepo[models.Song]) *Service {
	return &Service{
		repo:     repo,
		songRepo: songRepo,
	}
}

func (h *Service) AddSongToHistory(songYtId string, profileId uint) error {
	song, err := h.songRepo.GetByConds("yt_id = ?", songYtId)
	if err != nil {
		return err
	}

	return h.repo.Add(&models.History{
		ProfileId: profileId,
		SongId:    song[0].Id,
	})
}

func (h *Service) Get(profileId, page uint) ([]entities.Song, error) {
	gigaQuery := fmt.Sprintf(
		`SELECT yt_id, title, artist, thumbnail_url, duration, h.created_at
		FROM
			histories h JOIN songs
		ON
				songs.id = h.song_id
		WHERE h.profile_id = ?
		ORDER BY h.created_at DESC
		LIMIT %d,%d;`,
		(page-1)*20, page*20,
	)

	rows, err := h.repo.
		GetDB().
		Raw(gigaQuery, profileId).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := make([]entities.Song, 0)
	for rows.Next() {
		var song entities.Song
		var addedAt time.Time
		err = rows.Scan(&song.YtId, &song.Title, &song.Artist, &song.ThumbnailUrl, &song.Duration, &addedAt)
		if err != nil {
			continue
		}
		song.AddedAt = whenDidItHappen(addedAt)
		songs = append(songs, song)
	}

	songsFr := make([]entities.Song, 0)
	for i := 0; i < len(songs); i++ {
		playTimes := 1
		for i < len(songs)-1 && songs[i+1].YtId == songs[i].YtId {
			playTimes++
			i++
		}
		songs[i].AddedAt = fmt.Sprintf("Played %s - %s", times(playTimes), songs[i].AddedAt)
		songsFr = append(songsFr, songs[i])
	}

	return songsFr, nil
}

func whenDidItHappen(t time.Time) string {
	now := time.Now().UTC()
	switch {
	case t.Day() == now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "Today"
	case t.Day()+1 == now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "Yesterday"
	case t.Day()+5 < now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "Last week"
	case t.Day() == now.Day() && t.Month()+1 == now.Month() && t.Year() == now.Year():
		return "Last month"
	default:
		return fmt.Sprintf("%s %s %s", t.Format("January"), nth(t.Day()), t.Format("2006"))
	}
}

func nth(n int) string {
	switch {
	case n%10 == 1:
		return fmt.Sprintf("%dst", n)
	case n%10 == 2:
		return fmt.Sprintf("%dnd", n)
	case n%10 == 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
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

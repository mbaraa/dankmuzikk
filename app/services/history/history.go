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

	return songs, nil
}

func whenDidItHappen(t time.Time) string {
	now := time.Now().UTC()
	switch {
	case t.Day() == now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "today"
	case t.Day()+1 == now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "yesterday"
	case t.Day()+5 < now.Day() && t.Month() == now.Month() && t.Year() == now.Year():
		return "last week"
	case t.Day() == now.Day() && t.Month()+1 == now.Month() && t.Year() == now.Year():
		return "last month"
	default:
		return t.Format("2, January, 2006")
	}
}

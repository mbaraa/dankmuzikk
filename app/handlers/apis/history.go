package apis

import (
	"bytes"
	"dankmuzikk/entities"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/services/history"
	"dankmuzikk/views/components/song"
	"fmt"
	"net/http"
	"strconv"
)

type historyApi struct {
	service *history.Service
}

func NewHistoryApi(service *history.Service) *historyApi {
	return &historyApi{service}
}

func (h *historyApi) HandleGetMoreHistoryItems(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	page, err := strconv.Atoi(r.PathValue("page"))
	if err != nil {
		page = 2
	}
	if page <= 0 {
		page *= -1
	}
	recentPlays, err := h.service.Get(profileId, uint(page))
	if err != nil {
		log.Errorln(err)
	}
	if len(recentPlays) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	outBuf := bytes.NewBuffer([]byte{})
	for _, s := range recentPlays {
		song.Song(s, []string{"Played " + s.AddedAt}, nil, entities.Playlist{}).
			Render(r.Context(), outBuf)
	}

	outBuf.WriteString(fmt.Sprintf(`<div
			class="h-[10px] mb-[20px]"
			hx-get="/api/history/%d"
			hx-swap="outerHTML"
			hx-trigger="intersect"
			data-hx-revealed="true"
			data-loading-target="#history-loading"
			data-loading-class-remove="hidden"
			data-loading-path="/api/history/%d"
		></div>`, page+1, page+1))

	_, _ = w.Write(outBuf.Bytes())
}

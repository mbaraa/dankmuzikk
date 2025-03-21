package apis

import (
	"bytes"
	"dankmuzikk-web/entities"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/log"
	"dankmuzikk-web/services/history"
	"dankmuzikk-web/views/components/playlist"
	"dankmuzikk-web/views/components/song"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type historyApi struct {
	service *history.Service
}

func NewHistoryApi(service *history.Service) *historyApi {
	return &historyApi{service}
}

func (h *historyApi) HandleGetMoreHistoryItems(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
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

	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	recentPlays, err := h.service.Get(sessionToken.Value, uint(page))
	if err != nil {
		log.Errorln(err)
	}
	if len(recentPlays) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	outBuf := bytes.NewBuffer([]byte{})
	for idx, s := range recentPlays {
		song.Song(s, []string{s.AddedAt},
			[]templ.Component{
				playlist.PlaylistsPopup((idx+1)*page, s.YtId),
			},
			entities.Playlist{}).
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

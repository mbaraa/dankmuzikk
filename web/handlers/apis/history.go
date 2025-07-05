package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/playlist"
	"dankmuzikk-web/views/components/song"
	"dankmuzikk-web/views/components/status"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type historyApi struct {
	usecases *actions.Actions
}

func NewHistoryApi(usecases *actions.Actions) *historyApi {
	return &historyApi{
		usecases: usecases,
	}
}

func (h *historyApi) HandleGetMoreHistoryItems(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	page, err := strconv.Atoi(r.PathValue("page"))
	if err != nil {
		page = 2
	}
	if page <= 0 {
		page *= -1
	}

	recentPlays, err := h.usecases.GetHistory(actions.GetHistoryParams{
		ActionContext: ctx,
		PageIndex:     uint(page),
	})
	if err != nil {
		log.Errorln(err)
	}
	if len(recentPlays) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for idx, s := range recentPlays {
		song.Song(s, []string{s.AddedAt},
			[]templ.Component{
				playlist.PlaylistsPopup((idx+1)*page, s.PublicId),
			},
			actions.Playlist{}, "single").
			Render(r.Context(), w)
	}

	_, _ = w.Write(fmt.Appendf([]byte{}, `<div
	class="h-[10px] mb-[20px]"
	hx-get="/api/history/%d"
	hx-swap="outerHTML"
	hx-trigger="intersect"
	data-hx-revealed="true"
	data-loading-target="#history-loading"
	data-loading-class-remove="hidden"
	data-loading-path="/api/history/%d"></div>`,
		page+1, page+1))
}

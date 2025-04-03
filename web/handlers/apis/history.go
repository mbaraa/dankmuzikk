package apis

import (
	"bytes"
	"dankmuzikk-web/actions"
	"dankmuzikk-web/handlers/middlewares/auth"
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
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			BugsBunnyError("I'm not sure what you're trying to do here :)").
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

	recentPlays, err := h.usecases.GetHistory(sessionToken, uint(page))
	if err != nil {
		log.Errorln(err)
	}
	if len(recentPlays) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Warningln(recentPlays)

	outBuf := bytes.NewBuffer([]byte{})
	for idx, s := range recentPlays {
		song.Song(s, []string{s.AddedAt},
			[]templ.Component{
				playlist.PlaylistsPopup((idx+1)*page, s.PublicId),
			},
			actions.Playlist{}).
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

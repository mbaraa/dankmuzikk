package apis

import (
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

type libraryApi struct {
	usecases *actions.Actions
}

func NewLibraryApi(usecases *actions.Actions) *libraryApi {
	return &libraryApi{
		usecases: usecases,
	}
}

func (l *libraryApi) HandleAddSongToFavorites(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	songId := r.URL.Query().Get("id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := l.usecases.AddSongToFavorites(sessionToken, songId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song.RemoveFromFavoritesButton(songId).Render(r.Context(), w)
}

func (l *libraryApi) HandleRemoveSongFromFavorites(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			GenericError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	songId := r.URL.Query().Get("id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := l.usecases.RemoveSongFromFavorites(sessionToken, songId)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song.AddToFavoritesButton(songId).Render(r.Context(), w)
}

func (l *libraryApi) HandleGetMoreFavoritesItems(w http.ResponseWriter, r *http.Request) {
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

	payload, err := l.usecases.GetFavorites(sessionToken, uint(page))
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(payload.Songs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for idx, s := range payload.Songs {
		song.Song(s, []string{s.AddedAt},
			[]templ.Component{
				playlist.PlaylistsPopup((idx+1)*page, s.PublicId),
			},
			actions.Playlist{}, "favorites").
			Render(r.Context(), w)
	}

	w.Write(fmt.Appendf([]byte{}, `<div
			class="h-[10px] mb-[20px]"
			hx-get="/api/library/favorite/songs/%d"
			hx-swap="outerHTML"
			hx-trigger="intersect"
			data-hx-revealed="true"
			data-loading-target="#favorites-loading"
			data-loading-class-remove="hidden"
			data-loading-path="/api/library/favorite/songs/%d"></div>`,
		page+1, page+1))
}

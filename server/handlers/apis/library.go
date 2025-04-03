package apis

import (
	"dankmuzikk/actions"
	"dankmuzikk/log"
	"encoding/json"
	"net/http"
	"strconv"
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
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	songPublicId := r.URL.Query().Get("id")
	if songPublicId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("missing song's public id"))
		return
	}

	err = l.usecases.AddSongToFavorites(actions.AddSongToFavoritesParams{
		ActionContext: ctx,
		SongPublicId:  songPublicId,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}

func (l *libraryApi) HandleRemoveSongFromFavorites(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	songPublicId := r.URL.Query().Get("id")
	if songPublicId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("missing song's public id"))
		return
	}

	err = l.usecases.RemoveSongFromFavorites(actions.RemoveSongFromFavoritesParams{
		ActionContext: ctx,
		SongPublicId:  songPublicId,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}

func (l *libraryApi) HandleGetFavoriteSongs(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	page, err := strconv.Atoi(r.PathValue("page"))
	if err != nil || page <= 0 {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			handleErrorResponse(w, &ErrBadRequest{
				FieldName: "page",
			})
			return
		}
	}

	params := actions.GetFavoriteSongsParams{
		ActionContext: ctx,
		PageIndex:     uint(page),
	}

	recentPlays, err := l.usecases.GetFavoriteSongs(params)
	if err != nil {
		log.Error(err)
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(recentPlays)
}

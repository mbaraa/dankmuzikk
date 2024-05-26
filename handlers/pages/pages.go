package pages

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/history"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/youtube/download"
	"dankmuzikk/services/youtube/search"
	"dankmuzikk/views/layouts"
	"dankmuzikk/views/pages"
	"errors"
	"net/http"

	_ "github.com/a-h/templ"
)

const (
	notFoundMessage = "🤷‍♂️ I have no idea about what you requested!"
)

type pagesHandler struct {
	profileRepo      db.GetterRepo[models.Profile]
	playlistsService *playlists.Service
	jwtUtil          jwt.Manager[jwt.Json]
	ytSearch         search.Service
	downloadService  *download.Service
	historyService   *history.Service
}

func NewPagesHandler(
	profileRepo db.GetterRepo[models.Profile],
	playlistsService *playlists.Service,
	jwtUtil jwt.Manager[jwt.Json],
	ytSearch search.Service,
	downloadService *download.Service,
	historyService *history.Service,
) *pagesHandler {
	return &pagesHandler{
		profileRepo:      profileRepo,
		playlistsService: playlistsService,
		jwtUtil:          jwtUtil,
		ytSearch:         ytSearch,
		downloadService:  downloadService,
		historyService:   historyService,
	}
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	var recentPlays []entities.Song
	var err error
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if profileIdCorrect {
		recentPlays, err = p.historyService.Get(profileId)
		if err != nil {
			log.Errorln(err)
		}
	}

	if handlers.IsNoLayoutPage(r) {
		pages.Index(recentPlays).Render(r.Context(), w)
		return
	}
	layouts.Default(pages.Index(recentPlays)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if handlers.IsNoLayoutPage(r) {
		pages.About().Render(r.Context(), w)
		return
	}
	layouts.Default(pages.About()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw("Login", pages.Login()).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte(notFoundMessage))
		return
	}

	playlists, err := p.playlistsService.GetAll(profileId)
	if err != nil {
		playlists = make([]entities.Playlist, 0)
	}

	if handlers.IsNoLayoutPage(r) {
		pages.Playlists(playlists).Render(r.Context(), w)
		return
	}
	layouts.Default(pages.Playlists(playlists)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSinglePlaylistPage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if !profileIdCorrect {
		w.Write([]byte(notFoundMessage))
		return
	}

	playlistPubId := r.PathValue("playlist_id")
	if playlistPubId == "" {
		w.Write([]byte(notFoundMessage))
		return
	}

	playlist, permission, err := p.playlistsService.Get(playlistPubId, profileId)
	switch {
	case errors.Is(err, playlists.ErrUnauthorizedToSeePlaylist):
		log.Errorln(err)
		w.Write([]byte(notFoundMessage))
		return
	case err != nil:
		if playlist.Title == "" {
			log.Errorln(err)
			w.Write([]byte(notFoundMessage))
			return
		}
	}
	ctx := context.WithValue(r.Context(), handlers.PlaylistPermission, permission)

	if handlers.IsNoLayoutPage(r) {
		pages.Playlist(playlist).Render(ctx, w)
		return
	}
	layouts.Default(pages.Playlist(playlist)).Render(ctx, w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	layouts.Default(pages.Privacy()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if !profileIdCorrect {
		if handlers.IsNoLayoutPage(r) {
			w.Header().Set("HX-Redirect", "/")
		} else {
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		}
		return
	}
	// error is ignored, because the id was checked in the AuthHandler
	dbProfile, _ := p.profileRepo.Get(profileId)
	profile := entities.Profile{
		Name:     dbProfile.Name,
		PfpLink:  dbProfile.PfpLink,
		Username: dbProfile.Username,
	}
	if handlers.IsNoLayoutPage(r) {
		pages.Profile(profile).Render(r.Context(), w)
		return
	}
	layouts.Default(pages.Profile(profile)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSearchResultsPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results, err := p.ytSearch.Search(query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		log.Errorln(err)
		return
	}

	if len(results) != 0 {
		// TODO: move this call out of here
		log.Info("downloading songs' meta data from search")
		_ = p.downloadService.DownloadYoutubeSongsMetadata(results)
	}
	var songsInPlaylists map[string]bool
	var playlists []entities.Playlist
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if profileIdCorrect {
		playlists, songsInPlaylists, _ = p.playlistsService.GetAllMappedForAddPopover(profileId)
	}

	if handlers.IsNoLayoutPage(r) {
		pages.SearchResults(results, playlists, songsInPlaylists).Render(r.Context(), w)
		return
	}
	layouts.Default(pages.SearchResults(results, playlists, songsInPlaylists)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw("Sign up", pages.Signup()).Render(r.Context(), w)
}

package pages

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/handlers"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/youtube/search"
	"dankmuzikk/views/pages"
	"net/http"
	"strings"

	_ "github.com/a-h/templ"
)

const (
	notFoundMessage = "🤷‍♂️ I have no idea about what you requested!"
)

type pagesHandler struct {
	profileRepo      db.GetterRepo[models.Profile]
	playlistsService *playlists.Service
	jwtUtil          jwt.Manager[any]
}

func NewPagesHandler(
	profileRepo db.GetterRepo[models.Profile],
	playlistsService *playlists.Service,
	jwtUtil jwt.Manager[any],
) *pagesHandler {
	return &pagesHandler{profileRepo, playlistsService, jwtUtil}
}

func (p *pagesHandler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if handlers.IsNoReloadPage(r) {
		pages.AboutNoReload().Render(context.Background(), w)
		return
	}
	pages.About(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	if handlers.IsNoReloadPage(r) {
		pages.IndexNoReload().Render(context.Background(), w)
		return
	}
	pages.Index(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	pages.Login(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
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

	if handlers.IsNoReloadPage(r) {
		pages.PlaylistsNoReload(playlists).Render(context.Background(), w)
		return
	}
	pages.Playlists(p.isMobile(r), p.getTheme(r), playlists).Render(context.Background(), w)
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

	playlist, err := p.playlistsService.Get(playlistPubId, profileId)
	if err != nil {
		w.Write([]byte(notFoundMessage))
		return
	}
	_ = playlist

	if handlers.IsNoReloadPage(r) {
		// pages.PlaylistsNoReload(playlists).Render(context.Background(), w)
		return
	}
	// pages.Playlists(p.isMobile(r), p.getTheme(r), playlists).Render(context.Background(), w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	pages.Privacy(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
	if !profileIdCorrect {
		if handlers.IsNoReloadPage(r) {
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
	if handlers.IsNoReloadPage(r) {
		pages.ProfileNoReload(profile).Render(context.Background(), w)
		return
	}
	pages.Profile(p.isMobile(r), p.getTheme(r), profile).Render(context.Background(), w)
}

func (p *pagesHandler) HandleSearchResultsPage(ytSearch search.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		results, err := ytSearch.Search(query)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			log.Errorln(err)
			return
		}

		songs := make([]entities.Song, len(results))
		for i, result := range results {
			songs[i] = entities.Song{
				YtId:         result.Id,
				Title:        result.Title,
				Artist:       result.ChannelTitle,
				ThumbnailUrl: result.ThumbnailUrl,
				Duration:     result.Duration,
			}
		}

		var songsInPlaylists map[string]string
		var playlists []entities.Playlist
		profileId, profileIdCorrect := r.Context().Value(handlers.ProfileIdKey).(uint)
		if profileIdCorrect {
			log.Info("downloading songs from search")
			playlists, songsInPlaylists, _ = p.playlistsService.GetAllMappedForAddPopover(songs, profileId)
		}

		if handlers.IsNoReloadPage(r) {
			pages.SearchResultsNoReload(results, playlists, songsInPlaylists).Render(context.Background(), w)
			return
		}
		pages.SearchResults(p.isMobile(r), p.getTheme(r), results, playlists, songsInPlaylists).Render(context.Background(), w)
	}
}

func (p *pagesHandler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	pages.Signup(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}

func (p *pagesHandler) getTheme(r *http.Request) string {
	themeCookie, err := r.Cookie(handlers.ThemeName)
	if err != nil || themeCookie == nil || themeCookie.Value == "" {
		return "default"
	}
	switch themeCookie.Value {
	case "black":
		return "black"
	case "default":
		fallthrough
	default:
		return "default"
	}
}

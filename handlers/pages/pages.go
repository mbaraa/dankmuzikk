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
	"dankmuzikk/services/youtube/search"
	"dankmuzikk/views/pages"
	"net/http"
	"strings"

	_ "github.com/a-h/templ"
)

type pagesHandler struct {
	profileRepo db.GetterRepo[models.Profile]
	jwtUtil     jwt.Manager[any]
}

func NewPagesHandler(
	profileRepo db.GetterRepo[models.Profile],
	jwtUtil jwt.Manager[any],
) *pagesHandler {
	return &pagesHandler{profileRepo, jwtUtil}
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
	if handlers.IsNoReloadPage(r) {
		pages.PlaylistsNoReload().Render(context.Background(), w)
		return
	}
	pages.Playlists(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
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
		if handlers.IsNoReloadPage(r) {
			pages.SearchResultsNoReload(results).Render(context.Background(), w)
			return
		}
		pages.SearchResults(p.isMobile(r), p.getTheme(r), results).Render(context.Background(), w)
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

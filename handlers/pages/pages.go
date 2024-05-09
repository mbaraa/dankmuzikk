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
	"slices"
	"strings"

	_ "github.com/a-h/templ"
)

var noAuthPaths = []string{"/login", "/signup"}

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

func (p *pagesHandler) Handler(hand http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hand(w, r)
	}
}

func (p *pagesHandler) AuthHandler(hand http.HandlerFunc) http.HandlerFunc {
	return p.Handler(func(w http.ResponseWriter, r *http.Request) {
		htmxRedirect := p.isNoReload(r)
		authed := p.isAuthed(r, p.jwtUtil)

		switch {
		case authed && slices.Contains(noAuthPaths, r.URL.Path):
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		case !authed && slices.Contains(noAuthPaths, r.URL.Path):
			hand(w, r)
		case !authed && htmxRedirect:
			w.Header().Set("HX-Redirect", "/login")
		case !authed && !htmxRedirect:
			http.Redirect(w, r, config.Env().Hostname+"/login", http.StatusTemporaryRedirect)
		default:
			hand(w, r)
		}

	})
}

func (p *pagesHandler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if p.isNoReload(r) {
		pages.AboutNoReload().Render(context.Background(), w)
		return
	}
	pages.About(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	if p.isNoReload(r) {
		pages.IndexNoReload().Render(context.Background(), w)
		return
	}
	pages.Index(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	pages.Login(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	if p.isNoReload(r) {
		pages.PlaylistsNoReload().Render(context.Background(), w)
		return
	}
	pages.Playlists(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	pages.Privacy(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	tokenPayload := p.getRequestSessionTokenPayload(r)
	dbProfile, err := p.profileRepo.GetByConds("username = ?", tokenPayload["username"].(string))
	if err != nil {
		if p.isNoReload(r) {
			w.Header().Set("HX-Redirect", "/")
		} else {
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		}
		return
	}

	profile := entities.Profile{
		Name:     dbProfile[0].Name,
		PfpLink:  dbProfile[0].PfpLink,
		Username: dbProfile[0].Username,
	}

	if p.isNoReload(r) {
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
		if p.isNoReload(r) {
			pages.SearchResultsNoReload(results).Render(context.Background(), w)
			return
		}
		pages.SearchResults(p.isMobile(r), p.getTheme(r), results).Render(context.Background(), w)
	}
}

func (p *pagesHandler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	pages.Signup(p.isMobile(r), p.getTheme(r)).Render(context.Background(), w)
}

func (p *pagesHandler) isAuthed(r *http.Request, jwtUtil jwt.Manager[any]) bool {
	sessionToken, err := r.Cookie(handlers.SessionTokenKey)
	if err != nil {
		return false
	}
	err = jwtUtil.Validate(sessionToken.Value, jwt.SessionToken)
	if err != nil {
		return false
	}

	return true
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

func (p *pagesHandler) isNoReload(r *http.Request) bool {
	noReload, exists := r.URL.Query()["no_reload"]
	return exists && noReload[0] == "true"
}

func (p *pagesHandler) getRequestSessionTokenPayload(r *http.Request) map[string]any {
	// errors are ignored here becasue this method is used in pages that are wrapped with AuthHandler
	sessionToken, _ := r.Cookie(handlers.SessionTokenKey)
	token, _ := p.jwtUtil.Decode(sessionToken.Value, jwt.SessionToken)
	return token.Payload.(map[string]any)
}

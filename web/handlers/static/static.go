package static

import (
	"dankmuzikk-web/config"
	"dankmuzikk-web/static"
	"net/http"
	"strings"

	"github.com/tdewolff/minify/v2"
)

func HandleRobots(w http.ResponseWriter, r *http.Request) {
	robotsFile, _ := static.FS().ReadFile("robots.txt")
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write(robotsFile)
}

func HandleSitemap(w http.ResponseWriter, r *http.Request) {
	robotsFile, _ := static.FS().ReadFile("sitemap.xml")
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write(robotsFile)
}

func HandleFavicon(w http.ResponseWriter, r *http.Request) {
	faviconFile, _ := static.FS().ReadFile("favicon.ico")
	w.Header().Set("Content-Type", "image/x-icon")
	_, _ = w.Write(faviconFile)
}

func AssetsHandler(minifyer *minify.M) http.Handler {
	handler := http.NewServeMux()

	switch config.Env().GoEnv {
	case "prod":
		handler.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=7200, stale-while-revalidate=5")

			minifyer.Middleware(http.FileServer(http.FS(static.FS()))).
				ServeHTTP(w, r)
		}))
	case "beta":
		handler.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Env().GoEnv != "dev" {
				switch {
				case strings.HasPrefix(r.URL.Path, "/static/js"):
					w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=5")
				case strings.HasPrefix(r.URL.Path, "/static/css"):
					w.Header().Set("Cache-Control", "public, max-age=600, stale-while-revalidate=5")
				default:
					w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=5")
				}
			}

			http.FileServer(http.FS(static.FS())).
				ServeHTTP(w, r)
		}))
	default:
		handler.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.FS(static.FS())).
				ServeHTTP(w, r)
		}))
	}

	return handler
}

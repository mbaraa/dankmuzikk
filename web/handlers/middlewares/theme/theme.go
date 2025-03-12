package theme

import (
	"context"
	"net/http"
)

const ThemeKey = "theme-name"

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ThemeKey, getTheme(r))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTheme(r *http.Request) string {
	themeCookie, err := r.Cookie(ThemeKey)
	if err != nil || themeCookie == nil || themeCookie.Value == "" {
		return "black"
	}
	switch themeCookie.Value {
	case "dank":
		return "dank"
	case "white":
		return "white"
	case "black":
		fallthrough
	default:
		return "black"
	}
}

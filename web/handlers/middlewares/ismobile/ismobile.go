package ismobile

import (
	"context"
	"net/http"
	"strings"
)

const IsMobileKey = "is-mobile"

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), IsMobileKey, isMobile(r))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isMobile(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "mobile")
}

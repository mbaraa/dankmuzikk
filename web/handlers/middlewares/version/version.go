package version

import (
	"context"
	"net/http"
)

const VersionKey = "version"

func Handler(version string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), VersionKey, version)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

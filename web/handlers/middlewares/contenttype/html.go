package contenttype

import "net/http"

func Html(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h(w, r)
	}
}

// IsNoLayoutPage checks if the requested page requires a no reload or not.
func IsNoLayoutPage(r *http.Request) bool {
	noReload, exists := r.URL.Query()["no_layout"]
	return exists && noReload[0] == "true"
}

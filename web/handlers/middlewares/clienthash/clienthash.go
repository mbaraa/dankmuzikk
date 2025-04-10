package clienthash

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

const (
	ClientHashKey = "client-hash"
)

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := fmt.Sprintf(
			"%s:%s",
			r.Header.Get("X-Forwarded-For"),
			r.Header.Get("User-Agent"),
		)
		hasher := md5.New()
		hasher.Write([]byte(clientId))
		clientHash := hex.EncodeToString(hasher.Sum(nil))

		ctx := context.WithValue(r.Context(), ClientHashKey, clientHash)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

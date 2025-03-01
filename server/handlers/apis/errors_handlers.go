package apis

import (
	"dankmuzikk/app"
	"encoding/json"
	"net/http"
	"strings"
)

type errorResponse struct {
	ErrorId   string         `json:"error_id"`
	ExtraData map[string]any `json:"extra_data,omitempty"`
}

func handleErrorResponse(w http.ResponseWriter, err error) {
	if dankError, ok := err.(app.DankError); ok {
		if dankError.ExposeToClients() {
			w.WriteHeader(dankError.ClientStatusCode())
			_ = json.NewEncoder(w).Encode(errorResponse{
				ErrorId:   strings.ToLower(dankError.Error()),
				ExtraData: dankError.ExtraData(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(errorResponse{
		ErrorId: "internal-server-error",
	})
}

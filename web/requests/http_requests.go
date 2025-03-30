package requests

import (
	"bytes"
	"dankmuzikk-web/config"
	"dankmuzikk-web/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var r *requester

func init() {
	r = &requester{
		mu: sync.Mutex{},
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func getRequestUrl(path string) string {
	return fmt.Sprintf("%s%s", config.Env().ServerAddress, path)
}

type requester struct {
	mu         sync.Mutex
	httpClient *http.Client
}

func (r *requester) client() *http.Client {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.httpClient
}

type errorResponse struct {
	ErrorId   string         `json:"error_id"`
	ExtraData map[string]any `json:"extra_data,omitempty"`
}

type makeRequestConfig[T any] struct {
	method      string
	endpoint    string
	headers     map[string]string
	queryParams map[string]string
	body        T
}

func makeRequest[RequestBody any, ResponseBody any](conf makeRequestConfig[RequestBody]) (ResponseBody, error) {
	requestUrl := getRequestUrl(conf.endpoint)

	var respBody ResponseBody
	var bodyReader io.Reader = http.NoBody

	reqBodyType := reflect.TypeOf(conf.body)
	if reqBodyType != nil && reqBodyType.Kind() != reflect.Interface {
		bodyReaderLoc := bytes.NewBuffer(nil)
		err := json.NewEncoder(bodyReaderLoc).Encode(conf.body)
		if err != nil {
			return respBody, err
		}
		bodyReader = bodyReaderLoc
	} else {
		bodyReader = http.NoBody
	}

	req, err := http.NewRequest(conf.method, requestUrl, bodyReader)
	if err != nil {
		return respBody, err
	}

	q := req.URL.Query()
	for key, value := range conf.queryParams {
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range conf.headers {
		req.Header.Set(key, value)
	}

	resp, err := r.client().Do(req)
	if err != nil {
		return respBody, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return respBody, err
		}

		_ = resp.Body.Close()

		return respBody, mapError(errResp.ErrorId)
	}

	respBodyType := reflect.TypeOf(respBody)
	if respBodyType != nil && respBodyType.Kind() != reflect.Interface {
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return respBody, err
		}

		_ = resp.Body.Close()
	}

	return respBody, nil
}

func mapError(errorId string) error {
	switch errorId {
	case "invalid-token":
		return errors.ErrInvalidToken
	case "expired-token":
		return errors.ErrExpiredToken
	case "account-not-found":
		return errors.ErrAccountNotFound
	case "profile-not-found":
		return errors.ErrProfileNotFound
	case "song-not-found":
		return errors.ErrSongNotFound
	case "playlist-not-found":
		return errors.ErrPlaylistNotFound
	case "account-exists":
		return errors.ErrAccountExists
	case "profile-exists":
		return errors.ErrProfileExists
	case "song-exists":
		return errors.ErrSongExists
	case "playlist-exists":
		return errors.ErrPlaylistExists
	case "different-login-method-used":
		return errors.ErrDifferentLoginMethodUsed
	case "verification-code-expired":
		return errors.ErrVerificationCodeExpired
	case "invalid-verification-code":
		return errors.ErrInvalidVerificationCode
	case "non-owner-cant-delete-playlists":
		return errors.ErrNonOwnerCantDeletePlaylists
	case "non-owner-cant-change-playlist-visibility":
		return errors.ErrNonOwnerCantChangePlaylistVisibility
	case "unauthorized-to-see-playlist":
		return errors.ErrUnauthorizedToSeePlaylist
	case "user-has-already-voted":
		return errors.ErrUserHasAlreadyVoted
	case "not-enough-permission-to-add-song-to-playlist":
		return errors.ErrNotEnoughPermissionToAddSongToPlaylist
	default:
		return errors.ErrSomethingWentWrong
	}
}

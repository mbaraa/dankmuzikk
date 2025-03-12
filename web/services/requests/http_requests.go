package requests

import (
	"bytes"
	"dankmuzikk-web/config"
	"encoding/json"
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

func GetRequest[ResponseBody any](path string) (ResponseBody, error) {
	return getRequest[ResponseBody](path, map[string]string{})
}

func GetRequestAuthNoRespBody(path, token string) error {
	_, err := makeRequest[any, any](http.MethodGet, path, map[string]string{
		"Authorization": token,
	}, nil)
	return err
}

func GetRequestAuth[ResponseBody any](path, token string) (ResponseBody, error) {
	return getRequest[ResponseBody](path, map[string]string{
		"Authorization": token,
	})
}

func PostRequest[RequestBody any, ResponseBody any](path string, body RequestBody) (ResponseBody, error) {
	return postRequest[RequestBody, ResponseBody](path, body, map[string]string{})
}

func PostRequestAuth[RequestBody any, ResponseBody any](path, token string, body RequestBody) (ResponseBody, error) {
	return postRequest[RequestBody, ResponseBody](path, body, map[string]string{
		"Authorization": token,
	})
}

func PostRequestAuthNoBody[RequestBody any](path, token string, body RequestBody) error {
	_, err := makeRequest[RequestBody, any](http.MethodPost, path, map[string]string{}, body)
	return err
}

func PutRequest[RequestBody any, ResponseBody any](path string, body RequestBody) (ResponseBody, error) {
	return putRequest[RequestBody, ResponseBody](path, body, map[string]string{})
}

func PutRequestAuth[RequestBody any, ResponseBody any](path, token string, body RequestBody) (ResponseBody, error) {
	return putRequest[RequestBody, ResponseBody](path, body, map[string]string{
		"Authorization": token,
	})
}

func PutRequestAuthNoRespBody[RequestBody any](path, token string, body RequestBody) error {
	_, err := makeRequest[RequestBody, any](http.MethodPut, path, map[string]string{}, body)
	return err
}

func DeleteRequest(path string) error {
	return deleteRequest(path, map[string]string{})
}

func DeleteRequestAuth(path, token string) error {
	return deleteRequest(path, map[string]string{
		"Authorization": token,
	})
}

func getRequest[ResponseBody any](path string, headers map[string]string) (ResponseBody, error) {
	return makeRequest[any, ResponseBody](http.MethodGet, path, headers, nil)
}

func postRequest[RequestBody any, ResponseBody any](path string, body RequestBody, headers map[string]string) (ResponseBody, error) {
	return makeRequest[RequestBody, ResponseBody](http.MethodPost, path, headers, body)
}

func putRequest[RequestBody any, ResponseBody any](path string, body RequestBody, headers map[string]string) (ResponseBody, error) {
	return makeRequest[RequestBody, ResponseBody](http.MethodPut, path, headers, body)
}

func deleteRequest(path string, headers map[string]string) error {
	_, err := makeRequest[any, any](http.MethodDelete, path, headers, nil)
	return err
}

func makeRequest[RequestBody any, ResponseBody any](method, path string, headers map[string]string, body RequestBody) (ResponseBody, error) {
	var respBody ResponseBody

	var bodyReader io.Reader = http.NoBody

	reqBodyType := reflect.TypeOf(body)
	if reqBodyType != nil && reqBodyType.Kind() != reflect.Interface {
		bodyReaderLoc := bytes.NewBuffer(nil)
		err := json.NewEncoder(bodyReaderLoc).Encode(body)
		if err != nil {
			return respBody, err
		}
		bodyReader = bodyReaderLoc
	} else {
		bodyReader = http.NoBody
	}

	req, err := http.NewRequest(method, config.GetRequestUrl(path), bodyReader)
	if err != nil {
		return respBody, err
	}

	for key, value := range headers {
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

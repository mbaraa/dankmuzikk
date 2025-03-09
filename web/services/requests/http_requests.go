package requests

import (
	"bytes"
	"dankmuzikk-web/config"
	"dankmuzikk-web/log"
	"encoding/json"
	"errors"
	"net/http"
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

func GetRequest[ResponseBody any](path string) (ResponseBody, error) {
	return getRequest[ResponseBody](path, map[string]string{})
}

func GetRequestAuthNoRespBody(path, token string) error {
	req, err := http.NewRequest(http.MethodGet, config.GetRequestUrl(path), http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := r.client().Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("got %d status, when requesting GET %s\n", resp.StatusCode, path)
		return errors.New("non 200 status")
	}

	return nil
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
	bodyBytes := bytes.NewBuffer(nil)
	err := json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, config.GetRequestUrl(path), bodyBytes)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := r.client().Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("got %d status, when requesting %s %s\n", resp.StatusCode, "POST", path)
		return errors.New("non 200 status")
	}

	return nil
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
	bodyBytes := bytes.NewBuffer(nil)
	err := json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, config.GetRequestUrl(path), bodyBytes)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := r.client().Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("got %d status, when requesting PUT %s\n", resp.StatusCode, path)
		return errors.New("non 200 status")
	}

	return nil
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
	var respBody ResponseBody

	req, err := http.NewRequest(http.MethodGet, config.GetRequestUrl(path), http.NoBody)
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
		log.Errorf("got %d status, when requesting GET %s\n", resp.StatusCode, path)
		return respBody, errors.New("non 200 status")
	}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return respBody, err
	}

	_ = resp.Body.Close()

	return respBody, nil
}

func postRequest[RequestBody any, ResponseBody any](path string, body RequestBody, headers map[string]string) (ResponseBody, error) {
	var respBody ResponseBody

	bodyBytes := bytes.NewBuffer(nil)
	err := json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
		return respBody, err
	}

	req, err := http.NewRequest(http.MethodPost, config.GetRequestUrl(path), bodyBytes)
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
		log.Errorf("got %d status, when requesting POST %s\n", resp.StatusCode, path)
		return respBody, errors.New("non 200 status")
	}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return respBody, err
	}

	_ = resp.Body.Close()

	return respBody, nil
}

func putRequest[RequestBody any, ResponseBody any](path string, body RequestBody, headers map[string]string) (ResponseBody, error) {
	var respBody ResponseBody

	bodyBytes := bytes.NewBuffer(nil)
	err := json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
		return respBody, err
	}

	req, err := http.NewRequest(http.MethodPut, config.GetRequestUrl(path), bodyBytes)
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
		log.Errorf("got %d status, when requesting PUT %s\n", resp.StatusCode, path)
		return respBody, errors.New("non 200 status")
	}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return respBody, err
	}

	_ = resp.Body.Close()

	return respBody, nil
}

func deleteRequest(path string, headers map[string]string) error {
	req, err := http.NewRequest(http.MethodDelete, config.GetRequestUrl(path), http.NoBody)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := r.client().Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("got %d status, when requesting DELETE %s\n", resp.StatusCode, path)
		return errors.New("non 200 status")
	}

	return nil
}

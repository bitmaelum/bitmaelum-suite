package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// OutputResponse is a generic output response
type OutputResponse struct {
	Error  bool   `json:"error,omitempty"`
	Status string `json:"status"`
}

// JSONOut outputs the given data structure to JSON
func JSONOut(w http.ResponseWriter, v interface{}) error {
	var data []byte
	var err error

	// Indent or not, depends on the x-pretty-json header set in the response. This is set by the prettyjson middleware
	if w.Header().Get("x-pretty-json") == "1" {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}
	w.Header().Del("x-pretty-json")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(StatusError("Malformed JSON: " + err.Error()))
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	return err
}

// ErrorOut outputs an error
func ErrorOut(w http.ResponseWriter, code int, msg string) {
	logrus.Debugf("Returning error (%d): %s", code, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(StatusError(msg))
	return
}

// StatusOk Return an ok status response
func StatusOk(status string) *OutputResponse {
	return &OutputResponse{
		Status: status,
	}
}

// StatusError Return an error status response
func StatusError(status string) *OutputResponse {
	return &OutputResponse{
		Error:  true,
		Status: status,
	}
}

// DecodeBody decodes a JSON body or write/return an error if an error occurs
func DecodeBody(w http.ResponseWriter, body io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(StatusError("Malformed JSON: " + err.Error()))
		return err
	}

	return nil
}

// GetAPIKey returns the api key stored in the request context. If not found, it will return a dummy key with no permissions
func GetAPIKey(req *http.Request) *apikey.KeyType {
	val := req.Context().Value(middleware.APIKeyContext("apikey"))
	if val == nil {
		return &apikey.KeyType{}
	}

	return val.(*apikey.KeyType)
}

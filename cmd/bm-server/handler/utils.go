package handler

import (
	"encoding/json"
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
	data, err := json.MarshalIndent(v, "", "  ")
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

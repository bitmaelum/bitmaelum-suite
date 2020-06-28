package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OutputResponse struct {
	Error  bool   `json:"error,omitempty"`
	Status string `json:"status"`
}

func JsonOut(w http.ResponseWriter, v interface{}) error {
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

func ErrorOut(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(StatusError(msg))
	return
}

// Return an ok status response
func StatusOk(status string) *OutputResponse {
	return &OutputResponse{
		Status: status,
	}
}

// Return an error status response
func StatusError(status string) *OutputResponse {
	return &OutputResponse{
		Error:  true,
		Status: status,
	}
}

// Return an error status response
func StatusErrorf(status string, args ...interface{}) *OutputResponse {
	return &OutputResponse{
		Error:  true,
		Status: fmt.Sprintf(status, args...),
	}
}

// Decode a JSON body or write/return an error
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

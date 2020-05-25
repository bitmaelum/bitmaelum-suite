package handler

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type OutputResponse struct {
    Error bool `json:"error"`
    Status string `json:"status"`
}

func StatusOk(status string) *OutputResponse {
    return &OutputResponse{
        Error: false,
        Status: status,
    }
}
func StatusError(status string) *OutputResponse {
    return &OutputResponse{
        Error: true,
        Status: status,
    }
}

func StatusErrorf(status string, args ...interface{}) *OutputResponse {
    return &OutputResponse{
        Error: true,
        Status: fmt.Sprintf(status, args...),
    }
}


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

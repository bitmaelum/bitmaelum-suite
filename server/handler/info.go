package handler

import (
    "encoding/json"
    "net/http"
)

// Information header on root /
func Info(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(StatusOk("info"))
}

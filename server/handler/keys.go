package handler

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/core/container"
    "net/http"
)

type OutputPublicKey struct {
    PublicKey string `json:"public_key"`
}

type InputPublicKey struct {
    PublicKey string `json:"public_key"`
}


func RetrieveKey(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    hash := vars["sha256"]

    as := container.GetAccountService()
    if ! as.AccountExists(hash) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        _ = json.NewEncoder(w).Encode(StatusError("public key not found"))
        return
    }

    ret := OutputPublicKey{
        PublicKey: as.GetPublicKey(hash),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(ret)
}



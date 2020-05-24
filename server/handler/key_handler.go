package handler

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/core/keys"
    "net/http"
    "strings"
)

type OutputPublicKey struct {
    PublicKey string `json:"public_key"`
}

type InputPublicKey struct {
    PublicKey string `json:"public_key"`
}


func RetrieveKey(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    hash := strings.ToLower(vars["sha256"])

    if ! keys.HasKey(hash) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        _ = json.NewEncoder(w).Encode(StatusError("public key not found"))
        return
    }

    ret := OutputPublicKey{
        PublicKey: keys.GetKey(hash),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(ret)
}

func DeleteKey(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    hash := strings.ToLower(vars["sha256"])

    if ! keys.HasKey(hash) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        _ = json.NewEncoder(w).Encode(StatusError("public key not found"))
        return
    }

    keys.RemoveKey(hash)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(StatusOk("public key has been deleted"))
}

func StoreKey(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    hash := strings.ToLower(vars["sha256"])

    decoder := json.NewDecoder(req.Body)

    var input InputPublicKey
    err := decoder.Decode(&input)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusError("Malformed JSON: " + err.Error()))
        return
    }

    keys.AddKey(hash, input.PublicKey)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(StatusOk("public key has been added"))
}


package handler

import (
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/container"
    "net/http"
)

// Get flags from message
func GetFlags(w http.ResponseWriter, req *http.Request) {
    addr := core.HashAddress(mux.Vars(req)["addr"])
    box := mux.Vars(req)["box"]
    id := mux.Vars(req)["id"]

    // Retrieve flags
    as := container.GetAccountService()
    flags, _ := as.GetFlags(addr, box, id)

    _ = JsonOut(w, flags)
}

// Set flags for message
func SetFlag(w http.ResponseWriter, req *http.Request) {
    addr := core.HashAddress(mux.Vars(req)["addr"])
    box := mux.Vars(req)["box"]
    id := mux.Vars(req)["id"]
    flag := mux.Vars(req)["flag"]

    as := container.GetAccountService()
    _ = as.SetFlag(addr, box, id, flag)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}

// Unset flags for message
func UnsetFlag(w http.ResponseWriter, req *http.Request) {
    addr := core.HashAddress(mux.Vars(req)["addr"])
    box := mux.Vars(req)["box"]
    id := mux.Vars(req)["id"]
    flag := mux.Vars(req)["flag"]

    as := container.GetAccountService()
    _ = as.UnsetFlag(addr, box, id, flag)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}


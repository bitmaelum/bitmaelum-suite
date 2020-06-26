package handler

import (
    "encoding/json"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/gorilla/mux"
    "net/http"
)

type OutputPublicKey struct {
    PublicKeys []string `json:"public_key"`
}

type InputPublicKey struct {
    PublicKey string `json:"public_key"`
}

// Retrieve key handler
func RetrieveKeys(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    addr := core.HashAddress(vars["addr"])

    // Check if account exists
    as := container.GetAccountService()
    if ! as.AccountExists(addr) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        _ = json.NewEncoder(w).Encode(StatusError("public key not found"))
        return
    }

    // Return public key
    ret := OutputPublicKey{
        PublicKeys: as.GetPublicKeys(addr),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(ret)
}

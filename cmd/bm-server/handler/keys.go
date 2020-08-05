package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gorilla/mux"
	"net/http"
)

type outputPublicKey struct {
	PublicKeys []string `json:"public_key"`
}

// type inputPublicKey struct {
// 	PublicKey string `json:"public_key"`
// }

// RetrieveKeys is the handler that will retrieve public keys directly from the mailserver
func RetrieveKeys(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	// Check if account exists
	as := container.GetAccountService()
	if !as.AccountExists(*haddr) {
		ErrorOut(w, http.StatusNotFound, "public key not found")
		return
	}

	// Return public key
	ret := outputPublicKey{
		PublicKeys: as.GetPublicKeys(*haddr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ret)
}

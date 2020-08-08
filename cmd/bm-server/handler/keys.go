package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gorilla/mux"
	"net/http"
)

// RetrieveKeys is the handler that will retrieve public keys directly from the mailserver
func RetrieveKeys(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	// Check if account exists
	ar := container.GetAccountRepo()
	if !ar.Exists(*haddr) {
		ErrorOut(w, http.StatusNotFound, "public keys not found")
		return
	}

	keys, err := ar.FetchKeys(*haddr)
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "public keys not found")
		return
	}

	// Return public keys
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"public_keys": keys,
	})
}

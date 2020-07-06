package handler

import (
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	"github.com/gorilla/mux"
	"net/http"
)

// GetFlags Get flags from message
func GetFlags(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		// @TODO: Return error
		return
	}
	box := mux.Vars(req)["box"]
	id := mux.Vars(req)["id"]

	// Retrieve flags
	as := container.GetAccountService()
	flags, _ := as.GetFlags(*haddr, box, id)

	_ = JSONOut(w, flags)
}

// SetFlag Set flags for message
func SetFlag(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		// @TODO: Return error
		return
	}
	box := mux.Vars(req)["box"]
	id := mux.Vars(req)["id"]
	flag := mux.Vars(req)["flag"]

	as := container.GetAccountService()
	_ = as.SetFlag(*haddr, box, id, flag)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// UnsetFlag Unset flags for message
func UnsetFlag(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		// @TODO: Return error
		return
	}
	box := mux.Vars(req)["box"]
	id := mux.Vars(req)["id"]
	flag := mux.Vars(req)["flag"]

	as := container.GetAccountService()
	_ = as.UnsetFlag(*haddr, box, id, flag)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

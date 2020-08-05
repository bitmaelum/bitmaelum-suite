package handler

import (
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// GetFlags Get flags from message
func GetFlags(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	id := mux.Vars(req)["id"]
	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	// Retrieve flags
	as := container.GetAccountService()
	flags, _ := as.GetFlags(*haddr, box, id)

	_ = JSONOut(w, flags)
}

// SetFlag Set flags for message
func SetFlag(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}
	id := mux.Vars(req)["id"]
	flag := mux.Vars(req)["flag"]
	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	as := container.GetAccountService()
	_ = as.SetFlag(*haddr, box, id, flag)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// UnsetFlag Unset flags for message
func UnsetFlag(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	id := mux.Vars(req)["id"]
	flag := mux.Vars(req)["flag"]
	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	as := container.GetAccountService()
	_ = as.UnsetFlag(*haddr, box, id, flag)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

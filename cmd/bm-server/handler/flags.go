package handler

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
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

	messageID := mux.Vars(req)["message"]
	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	// Retrieve flags
	ar := container.GetAccountRepo()
	flags, _ := ar.GetFlags(*haddr, box, messageID)

	_ = JSONOut(w, flags)
}

// SetFlag Set flags for message
func SetFlag(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}
	messageID := mux.Vars(req)["message"]
	flag := mux.Vars(req)["flag"]
	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	ar := container.GetAccountRepo()
	_ = ar.SetFlag(*haddr, box, messageID, flag)

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

	messageID := mux.Vars(req)["message"]
	flag := mux.Vars(req)["flag"]
	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	ar := container.GetAccountRepo()
	_ = ar.UnsetFlag(*haddr, box, messageID, flag)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

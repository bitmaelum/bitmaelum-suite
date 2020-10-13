package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/parse"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
)

type inputAPIKeyType struct {
	Permissions []string `json:"permissions"`
	Valid       string   `json:"valid,omitempty"`
	Desc        string   `json:"description,omitempty"`
}

// CreateAPIKey is a handler that will create a new API key (non-admin keys only)
func CreateAPIKey(w http.ResponseWriter, req *http.Request) {
	var input inputAPIKeyType
	err := DecodeBody(w, req.Body, &input)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	// Our custom parser allows (and defaults) to using days
	validDuration, err := parse.ValidDuration(input.Valid)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect valid duration")
		return
	}
	//
	err = parse.AccountPermissions(input.Permissions)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect permissions")
		return
	}

	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	newAPIKey := apikey.NewAccountKey(h, input.Permissions, validDuration, input.Desc)

	// Store API key into persistent storage
	repo := container.GetAPIKeyRepo()
	err = repo.Store(newAPIKey)
	if err != nil {
		msg := fmt.Sprintf("error while storing key: %s", err)
		ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"api_key": newAPIKey.ID,
	})
}

// ListAPIKeys returns a list of all keys for the given account
func ListAPIKeys(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	// Store API key into persistent storage
	repo := container.GetAPIKeyRepo()
	keys, err := repo.FetchByHash(h.String())
	if err != nil {
		msg := fmt.Sprintf("error while retrieving keys: %s", err)
		ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(keys)
}

// DeleteAPIKey will remove a key
func DeleteAPIKey(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.GetAPIKeyRepo()
	key, err := repo.Fetch(keyID)
	if err != nil || key.AddrHash.String() != h.String() {
		// Only allow deleting of keys that we own as account
		ErrorOut(w, http.StatusNotFound, "key not found")
		return
	}

	_ = repo.Remove(*key)

	// All is well
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetAPIKey will get a key
func GetAPIKeyDetails(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.GetAPIKeyRepo()
	key, err := repo.Fetch(keyID)
	if err != nil || key.AddrHash.String() != h.String() {
		ErrorOut(w, http.StatusNotFound, "key not found")
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(key)
}

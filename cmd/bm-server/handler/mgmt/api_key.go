package mgmt

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/parse"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

type inputAPIKeyType struct {
	Permissions []string `json:"permissions"`
	Valid       string   `json:"valid"`
	AddrHash    string   `json:"hash,omitempty"`
	Desc        string   `json:"description,omitempty"`
}

// NewAPIKey is a handler that will create a new API key (non-admin keys only)
func NewAPIKey(w http.ResponseWriter, req *http.Request) {
	var input inputAPIKeyType
	err := handler.DecodeBody(w, req.Body, &input)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	// Make sure we can only set api permissions for the account we have permission for.
	var h *hash.Hash
	if input.AddrHash != "" {
		tmp := hash.New(input.AddrHash)
		h = &tmp
	}
	key := handler.GetAPIKey(req)
	if !key.HasPermission(apikey.PermAPIKeys, h) {
		handler.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}


	// Our custom parser allows (and defaults) to using days
	validDuration, err := parse.ValidDuration(input.Valid)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect valid duration")
		return
	}

	err = parse.Permissions(input.Permissions)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect permissions")
		return
	}

	newAPIKey := apikey.NewKey(input.Permissions, validDuration, h, input.Desc)

	// Store API key into persistent storage
	repo := container.GetAPIKeyRepo()
	err = repo.Store(newAPIKey)
	if err != nil {
		msg := fmt.Sprintf("error while storing key: %s", err)
		handler.ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"api_key": newAPIKey.ID,
	})
}

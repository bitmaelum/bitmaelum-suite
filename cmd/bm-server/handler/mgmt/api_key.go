package mgmt

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"net/http"
)

type inputAPIKeyType struct {
	Permissions []string `json:"permissions"`
	Valid       string   `json:"valid"`
}

// NewAPIKey is a handler that will create a new API key (non-admin keys only)
func NewAPIKey(w http.ResponseWriter, req *http.Request) {
	key := handler.GetAPIKey(req)
	if !key.HasPermission(apikey.PermAPIKeys) {
		handler.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input inputAPIKeyType
	err := handler.DecodeBody(w, req.Body, &input)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	// Our custom parser allows (and defaults) to using days
	validDuration, err := apikey.ParseValidDuration(input.Valid)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect valid duration")
		return
	}

	err = apikey.ParsePermissions(input.Permissions)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect permissions")
		return
	}

	newAPIKey := apikey.NewKey(input.Permissions, validDuration)

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

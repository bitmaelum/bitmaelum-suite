package mgmt

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/invite"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

type inputInviteType struct {
	AddrHash string `json:"address"`
	Days     int    `json:"days"`
}

type jsonOut map[string]interface{}

// NewInvite handler will generate a new invite token for a given address
func NewInvite(w http.ResponseWriter, req *http.Request) {
	key := handler.GetAPIKey(req)
	if !key.HasPermission(apikey.PermGenerateInvites) {
		handler.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input inputInviteType
	err := handler.DecodeBody(w, req.Body, &input)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	addrHash, err := hash.NewFromHash(input.AddrHash)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	validUntil := time.Now().Add(time.Duration(input.Days) * 24 * time.Hour)
	token, err := invite.NewInviteToken(*addrHash, config.Routing.RoutingID, validUntil, config.Routing.PrivateKey)
	if err != nil {
		handler.ErrorOut(w, http.StatusInternalServerError, "cannot generate invite token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"hash":   addrHash.String(),
		"token":  token.String(),
		"expiry": validUntil.Unix(),
	})
}

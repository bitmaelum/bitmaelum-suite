package mgmt

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"net/http"
	"time"
)

type inputInviteType struct {
	Addr string `json:"address"`
	Days int    `json:"days"`
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

	addr, err := address.NewHashFromHash(input.Addr)
	if err != nil {
		handler.ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	inviteRepo := container.GetInviteRepo()
	token, err := inviteRepo.Get(*addr)
	if err == nil {
		msg := fmt.Sprintf("'%s' already allowed to register with token: %s\n", addr.String(), token)
		handler.ErrorOut(w, http.StatusConflict, msg)
		return
	}

	token, err = inviteRepo.Create(*addr, time.Duration(input.Days)*24*time.Hour)
	if err != nil {
		msg := fmt.Sprintf("error while inviting address: %s", err)
		handler.ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"invite_token": token,
	})
}

package handler

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/internal/config"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	pow "github.com/bitmaelum/bitmaelum-server/pkg/proofofwork"
	"net/http"
)

type inputCreateAccount struct {
	Addr        address.HashAddress `json:"address"`
	Token       string              `json:"token"`
	PublicKey   string              `json:"public_key"`
	ProofOfWork struct {
		Bits  int    `json:"bits"`
		Proof uint64 `json:"proof"`
	} `json:"proof_of_work"`
}

// CreateAccount will create a new account
func CreateAccount(w http.ResponseWriter, req *http.Request) {
	var input inputCreateAccount
	err := DecodeBody(w, req.Body, &input)
	if err != nil {
		return
	}

	// Check proof of work first
	if input.ProofOfWork.Bits < config.Server.Accounts.ProofOfWork {
		ErrorOut(w, http.StatusBadRequest, fmt.Sprintf("Proof of work must be at least %d bits", config.Server.Accounts.ProofOfWork))
		return
	}

	work := pow.New(input.ProofOfWork.Bits, []byte(input.Addr), input.ProofOfWork.Proof)
	if !work.Validate() {
		ErrorOut(w, http.StatusBadRequest, "incorrect proof of work")
		return
	}

	// Check if token exists for the given address
	is := container.GetInviteService()
	registeredToken, err := is.GetInvite(input.Addr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "token not found")
		return
	}
	if subtle.ConstantTimeCompare([]byte(registeredToken), []byte(input.Token)) != 1 {
		ErrorOut(w, http.StatusBadRequest, "incorrect token")
		return
	}

	// Check if account exists
	as := container.GetAccountService()
	if as.AccountExists(input.Addr) {
		ErrorOut(w, http.StatusBadRequest, "account already exists")
		return
	}

	// All clear. Create account
	err = as.CreateAccount(input.Addr, input.PublicKey)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Done with the invite, let's remove
	_ = is.RemoveInvite(input.Addr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(StatusOk("BitMaelum account has been successfully created."))
}

// RetrieveAccount retrieves account information
func RetrieveAccount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("this is your account"))
	return
}

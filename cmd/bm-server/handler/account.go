package handler

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type inputCreateAccount struct {
	Addr        address.HashAddress `json:"address"`
	Token       string              `json:"token"`
	PublicKey   bmcrypto.PubKey     `json:"public_key"`
	ProofOfWork pow.ProofOfWork     `json:"proof_of_work"`
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

	work := pow.New(input.ProofOfWork.Bits, input.Addr.String(), input.ProofOfWork.Proof)
	if !work.IsValid() {
		ErrorOut(w, http.StatusBadRequest, "incorrect proof of work")
		return
	}

	// Check if token exists for the given address
	inviteRepo := container.GetInviteRepo()
	registeredToken, err := inviteRepo.Get(input.Addr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "token not found")
		return
	}
	if subtle.ConstantTimeCompare([]byte(registeredToken), []byte(input.Token)) != 1 {
		ErrorOut(w, http.StatusBadRequest, "incorrect token")
		return
	}

	// Check if account exists
	ar := container.GetAccountRepo()
	if ar.Exists(input.Addr) {
		ErrorOut(w, http.StatusBadRequest, "account already exists")
		return
	}

	// All clear. Create account
	err = ar.Create(input.Addr, input.PublicKey)
	if err != nil {
		logrus.Error(err)
		ErrorOut(w, http.StatusInternalServerError, "cannot create account")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(StatusOk("BitMaelum account has been successfully created."))
}

// RetrieveOrganisation is the handler that will retrieve organisation settings
func RetrieveOrganisation(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	// Check if account exists
	ar := container.GetAccountRepo()
	if !ar.Exists(*haddr) {
		ErrorOut(w, http.StatusNotFound, "address not found")
		return
	}

	settings, err := ar.FetchOrganisationSettings(*haddr)
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "organisation settings not found")
		return
	}

	// Return public keys
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(settings)
}

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

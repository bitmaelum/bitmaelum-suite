package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/invite"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type inputCreateAccount struct {
	Addr        address.Hash    `json:"address"`
	UserHash    string          `json:"user_hash"`
	OrgHash     string          `json:"org_hash"`
	Token       string          `json:"token"`
	PublicKey   bmcrypto.PubKey `json:"public_key"`
	ProofOfWork pow.ProofOfWork `json:"proof_of_work"`
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

	// Check if the user+org matches our actual hash address
	if !input.Addr.Verify(address.NewHash(input.UserHash), address.NewHash(input.OrgHash)) {
		ErrorOut(w, http.StatusBadRequest, "cant verify the address hashes")
		return
	}

	// Check if we need to verify against the mailserver key, or the organisation key
	var pubKey bmcrypto.PubKey = config.Server.Routing.PublicKey
	if input.OrgHash != "" {
		r := container.GetResolveService()
		oh, err := address.HashFromString(input.OrgHash)
		if err != nil {
			ErrorOut(w, http.StatusBadRequest, "incorrect org hash")
			return
		}
		oi, err := r.ResolveOrganisation(*oh)
		if err != nil {
			ErrorOut(w, http.StatusBadRequest, "cannot find organisation")
			return
		}
		pubKey = oi.PublicKey
	}

	// Verify token
	it, err := invite.ParseInviteToken(input.Token)
	if err != nil || !it.Verify(config.Server.Routing.RoutingID, pubKey) {
		ErrorOut(w, http.StatusBadRequest, "cannot validate token")
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
	haddr, err := address.HashFromString(mux.Vars(req)["addr"])
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
	haddr, err := address.HashFromString(mux.Vars(req)["addr"])
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

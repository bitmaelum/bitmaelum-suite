// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/signature"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type inputCreateAccount struct {
	Addr        hash.Hash       `json:"address"`
	UserHash    string          `json:"user_hash"`
	OrgHash     string          `json:"org_hash"`
	Token       string          `json:"token"`
	PublicKey   bmcrypto.PubKey `json:"public_key"`
	ProofOfWork pow.ProofOfWork `json:"proof_of_work"`
}

// CreateAccount will create a new account
func CreateAccount(w http.ResponseWriter, req *http.Request) {
	var input inputCreateAccount
	err := httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		return
	}

	// Get required number of bits from the resolver
	resolver := container.Instance.GetResolveService()
	cfg := resolver.GetConfig()

	// Check proof of work first
	if input.ProofOfWork.Bits < cfg.ProofOfWork.Address {
		httputils.ErrorOut(w, http.StatusBadRequest, fmt.Sprintf("Proof of work must be at least %d bits", cfg.ProofOfWork.Address))
		return
	}

	work := pow.New(input.ProofOfWork.Bits, input.Addr.String(), input.ProofOfWork.Proof)
	if !work.IsValid() {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect proof of work")
		return
	}

	// Check if the user+org matches our actual hash address
	userHash, err := hash.NewFromHash(input.UserHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "invalid body: user_hash")
		return
	}
	orgHash, err := hash.NewFromHash(input.OrgHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "invalid body: org_hash")
		return
	}
	if !input.Addr.Verify(*userHash, *orgHash) {
		httputils.ErrorOut(w, http.StatusBadRequest, "cant verify the address hashes")
		return
	}

	// Check if we need to verify against the mail server key, or the organisation key
	var pubKey = &config.Routing.KeyPair.PubKey
	if !orgHash.IsEmpty() {
		pubKey, err = checkOrganisation(input, w)
		if err != nil {
			return
		}
	}

	// Verify token
	it, err := signature.ParseInviteToken(input.Token)
	if err != nil || !it.Verify(config.Routing.RoutingID, *pubKey) {
		httputils.ErrorOut(w, http.StatusBadRequest, "cannot validate token")
		return
	}

	// Check if account exists
	ar := container.Instance.GetAccountRepo()
	if ar.Exists(input.Addr) {
		httputils.ErrorOut(w, http.StatusBadRequest, "account already exists")
		return
	}

	// All clear. Create account
	err = ar.Create(input.Addr, input.PublicKey)
	if err != nil {
		logrus.Error(err)
		httputils.ErrorOut(w, http.StatusInternalServerError, "cannot create account")
		return
	}

	_ = httputils.JSONOut(w, http.StatusCreated, httputils.StatusOk("BitMaelum account has been successfully created."))
}

func checkOrganisation(input inputCreateAccount, w http.ResponseWriter) (*bmcrypto.PubKey, error) {
	r := container.Instance.GetResolveService()
	orgHash, err := hash.NewFromHash(input.OrgHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect org hash")
		return nil, err
	}
	orgInfo, err := r.ResolveOrganisation(*orgHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "cannot find organisation")
		return nil, err
	}

	// Check if the organisation is whitelisted on our server
	allowed := false
	for _, org := range config.Server.Organisations {
		if hash.New(org).String() == orgHash.String() {
			allowed = true
		}
	}
	if !allowed {
		httputils.ErrorOut(w, http.StatusBadRequest, "organisation not allowed to register on this server")
		return nil, errors.New("org not alllowed")
	}

	return &orgInfo.PublicKey, nil
}

// RetrieveKeys is the handler that will retrieve public keys directly from the mail server
func RetrieveKeys(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	// Check if account exists
	ar := container.Instance.GetAccountRepo()
	if !ar.Exists(*haddr) {
		httputils.ErrorOut(w, http.StatusNotFound, "public keys not found")
		return
	}

	keys, err := ar.FetchKeys(*haddr)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, "public keys not found")
		return
	}

	// Return public keys
	_ = httputils.JSONOut(w, http.StatusOK, jsonOut{
		"public_keys": keys,
	})
}

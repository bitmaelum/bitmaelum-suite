package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/internal/message"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	proofOfWork        string = "proof_of_work"
	bodyAccept         string = "body_accept"
	bitsForProofOfWork int    = 22
)

// ProofOfWorkType structure
type ProofOfWorkType struct {
	Bits    int    `json:"bits"`
	Nonce   string `json:"nonce"`
	Path    string `json:"path"`
	Timeout string `json:"timeout"`
}

// BodyAcceptType structure
type BodyAcceptType struct {
	Path    string `json:"path"`
	Timeout string `json:"timeout"`
}

// OutputHeaderType structure
type OutputHeaderType struct {
	Error       bool             `json:"error"`
	Status      string           `json:"status"`
	Description string           `json:"description"`
	ProofOfWork *ProofOfWorkType `json:"proof_of_work,omitempty"`
	BodyAccept  *BodyAcceptType  `json:"body_accept,omitempty"`
}

// PostMessageHeader Handler when a message header is posted
func PostMessageHeader(w http.ResponseWriter, req *http.Request) {
	is := container.GetIncomingService()

	// Generate checksum for header message
	body, err := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	checksum := sha256.Sum256(body)

	// Decode JSON
	decoder := json.NewDecoder(req.Body)
	var input message.Header
	err = decoder.Decode(&input)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// @TODO: Validate incoming header
	//err = message.ValidateHeader(input)
	//if err != nil {
	//    sendBadRequest(w, err)
	//    return
	//}

	// Check if we need proof of work.
	if needsProofOfWork(input) {
		// Generate proof-of-work data
		path, nonce, err := is.GeneratePowResponsePath(input.From.Addr, bitsForProofOfWork, checksum[:])
		if err != nil {
			ErrorOut(w, http.StatusBadRequest, err.Error())
			return
		}

		// Allow 30 minutes for proof-of-work
		to := time.Now()
		to.Add(time.Minute * 30)

		pow := &ProofOfWorkType{
			Bits:    bitsForProofOfWork,
			Nonce:   nonce,
			Path:    "/incoming/" + path,
			Timeout: to.Format(time.RFC3339),
		}

		ret := OutputHeaderType{
			Error:       false,
			Status:      proofOfWork,
			Description: "A proof of work is needed before we will accept this message",
			ProofOfWork: pow,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ret)
		return
	}

	// No proof-of-work, generate accept path
	path, err := is.GenerateAcceptResponsePath(input.From.Addr, checksum[:])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Allow 30 minutes for incoming body message
	to := time.Now()
	to.Add(time.Minute * 30)

	ret := OutputHeaderType{
		Error:       false,
		Status:      bodyAccept,
		Description: "Accepting body for this header",
		BodyAccept: &BodyAcceptType{
			Path:    "/incoming/" + path,
			Timeout: to.Format(time.RFC3339),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ret)
}

// Defines if we need to do a proof-of-work based on the incoming header message
func needsProofOfWork(header message.Header) bool {
	// @TODO: We probably want to use different metrics to check if we need to do proof-of-work
	return rand.Intn(10) < 5
}

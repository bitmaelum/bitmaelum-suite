package handler

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/storage"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/gorilla/mux"
	mr "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTP Header that contains proof-of-work in <challenge>::<proof> format.
const proofOfWorkHeader = "x-bitmaelum-pow"

/*
 * Incoming is when a SERVER sends a message to another SERVER. This is an unauthenticated action so there is need for
 * proof-of-work, unless we have a valid subscription ID. The server only accepts messages for local accounts. There is
 * no relaying of any kind.
 *
 * Incoming works a bit different than upload. Incoming needs to request an upload first. If proof-of-work needs to be
 * completed, it will respond with a 412 message with a header { "data": "abcd1234...125", bits: 22 }. This will be the
 * proof of work that needs to be completed for THIS particular message. The request must be completed within a half
 * hour or the proof-of-work request will expire.
 * If validated correctly, the request will respond with a message uuid, where the user may upload the message to.
 *
 */

type jsonOut map[string]interface{}

// IncomingMessageRequest requests to upload a message. It might need a proof-of-work response and if ok, it will return
// the messageID which can be used for actual uploading.
func IncomingMessageRequest(w http.ResponseWriter, req *http.Request) {
	pow, err := getProofOfWorkFromHeader(req)
	if err != nil {
		// proof of work not found, generate a new proof of work
		pow, err = storage.NewProofOfWork()
		if err != nil {
			ErrorOut(w, http.StatusInternalServerError, "cannot generate challenge")
			return
		}
	} else {
		// Check if found proof of work actually is valid
		p := proofofwork.New(pow.Bits, []byte(pow.Challenge), pow.Proof)
		// if it was already valid, don't invalidate it because we posted something incorrect again
		pow.Valid = pow.Valid || (p.HasDoneWork() && p.IsValid())
	}

	// If we don't need proof of work, we are automatically valid
	if !needsProofOfWork(req) {
		pow.Valid = true
	}

	if pow.Valid {
		// We don't expire anymore since we can upload the message now. We remove the item explicitly once the
		// message upload has been completed
		pow.Expires = time.Unix(0, 0)
	}

	// Save proof of work
	powService := container.GetProofOfWorkService()
	err = powService.Store(pow)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing pow")
		return
	}

	// Check if the work was validated or we didn't need it
	if pow.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jsonOut{"MsgID": pow.MsgID})
		return
	}

	// Proof of work is not valid (or was not found). Return the challenge
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPreconditionFailed)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"challenge": pow.Challenge,
		"bits":      pow.Bits,
	})
}

// IncomingMessageHeader deals with uploading message headers
func IncomingMessageHeader(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	// Did we already upload the header?
	if message.IncomingPathExists(msgID, "header.json") {
		ErrorOut(w, http.StatusConflict, "header already uploaded")
		return
	}

	// Read header from request body
	header, err := readHeaderFromBody(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "invalid header")
		return
	}

	// Save request
	err = message.StoreHeader(msgID, header)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message header")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("header saved"))
	return
}

// IncomingMessageCatalog deals with uploading message catalogs
func IncomingMessageCatalog(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	// Did we already upload the header?
	if message.IncomingPathExists(msgID, "catalog") {
		ErrorOut(w, http.StatusConflict, "catalog already uploaded")
		return
	}

	err := message.StoreCatalog(msgID, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message header")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved catalog"))
	return
}

// IncomingMessageBlock deals with uploading message blocks
func IncomingMessageBlock(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]
	id := vars["id"]

	// Did we already upload the header?
	if message.IncomingPathExists(msgID, id) {
		ErrorOut(w, http.StatusConflict, "block already uploaded")
		return
	}

	err := message.StoreBlock(msgID, id, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message block")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved message block"))
	return
}

// CompleteIncoming is called whenever everything from a message has been uploaded and can be actually send
func CompleteIncoming(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	// @TODO: How do we know if all data is send over? Maybe we should add all files to the header so we can verify?
	if !message.IncomingPathExists(msgID, "header.json") || !message.IncomingPathExists(msgID, "catalog") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	// queue the message for processing
	processor.QueueIncomingMessage(msgID)
}

// DeleteIncoming is called whenever we want to completely remove a message by user request
func DeleteIncoming(w http.ResponseWriter, req *http.Request) {
	// Delete the message and contents
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	if !message.IncomingPathExists(msgID, "") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	err := message.RemoveMessage(message.SectionIncoming, msgID)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while deleting outgoing message")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("message removed"))
	return
}

// needsProofOfWork will check if the given data posted to the request allows for skipping the proof of work.
// This is done for instance when we posted a valid sendAddr-recipientAddr-subscriptionID tuple.
func needsProofOfWork(req *http.Request) bool {
	// @TODO: Let's figure out a better way to determinate if we need proof-of-work
	return mr.Intn(10) > 5
}

// Returns the proof of work found in the header, or err when not found / invalid
func getProofOfWorkFromHeader(req *http.Request) (*storage.ProofOfWork, error) {
	tmp := strings.SplitN(req.Header.Get(proofOfWorkHeader), "::", 2)
	if len(tmp) != 2 {
		return nil, errors.New("incorrect header format for proof of work")
	}
	challenge, proofStr := tmp[0], tmp[1]
	proof, err := strconv.Atoi(proofStr)
	if err != nil {
		return nil, errors.New("incorrect proof")
	}

	// Check if the POW exists in our storage
	powService := container.GetProofOfWorkService()
	pow, err := powService.Retrieve(challenge)
	if err != nil {
		return nil, errors.New("challenge not found")
	}

	// Set the proof of the storage
	pow.Proof = uint64(proof)
	return pow, nil
}

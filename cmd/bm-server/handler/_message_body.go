package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-server/cmd/bm-server/incoming"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	pow "github.com/bitmaelum/bitmaelum-server/pkg/proofofwork"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// PostMessageBody in progress
func PostMessageBody(w http.ResponseWriter, req *http.Request) {
	is := container.GetIncomingService()
	path := mux.Vars(req)["addr"]

	// Check if the path is actually an UUID
	_, err := uuid.Parse(path)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if this UUID path is an live incoming path (and not yet expired)
	info, err := is.GetIncomingPath(path)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	if info == nil {
		ErrorOut(w, http.StatusNotFound, "not found")
		return
	}

	// Handle either accept response path or proof-of-work response path
	switch info.Type {
	case incoming.Accept:
		handleAccept(w, req, info)
		return
	case incoming.ProofOfWork:
		handlePow(w, req, info)
		return
	}

	// Something else has happened. Unknown
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(StatusError("unknown incoming type for this request"))
}

func handlePow(w http.ResponseWriter, req *http.Request, info *incoming.InfoType) {
	is := container.GetIncomingService()
	decoder := json.NewDecoder(req.Body)

	// Decode JSON
	var input pow.ProofOfWork
	err := decoder.Decode(&input)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Make sure the proof-of-work is completed and valid
	work := pow.New(input.Bits, []byte(info.Nonce), input.Proof)
	if !work.Validate() {
		ErrorOut(w, http.StatusNotAcceptable, "incorrect proof of work")
		return
	}

	// We can generate an accept path
	path, err := is.GenerateAcceptResponsePath(info.Addr, info.Checksum)
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

func handleAccept(w http.ResponseWriter, req *http.Request, info *incoming.InfoType) {

}

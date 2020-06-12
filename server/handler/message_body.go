package handler

import (
    "encoding/json"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/container"
    "github.com/bitmaelum/bitmaelum-server/server/incoming"
    "net/http"
    "time"
)

//
func PostMessageBody(w http.ResponseWriter, req *http.Request) {
    is := container.GetIncomingService()
    path := mux.Vars(req)["addr"]

    // Check if the path is actually an UUID
    _, err := uuid.Parse(path)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    // Check if this UUID path is an live incoming path (and not yet expired)
    info, err := is.GetIncomingPath(path)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    if info == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        _ = json.NewEncoder(w).Encode(StatusError("not found"))
        return
    }

    // Handle either accept response path or proof-of-work response path
    switch info.Type {
    case incoming.ACCEPT:
        handleAccept(w, req, info)
        return
    case incoming.PROOF_OF_WORK:
        handlePow(w, req, info)
        return
    }

    // Something else has happened. Unknown
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    _ = json.NewEncoder(w).Encode(StatusError("unknown incoming type for this request"))
}

func handlePow(w http.ResponseWriter, req *http.Request, info *incoming.IncomingInfoType) {
    is := container.GetIncomingService()
    decoder := json.NewDecoder(req.Body)

    // Decode JSON
    var input core.ProofOfWork
    err := decoder.Decode(&input)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    // Make sure the proof-of-work is completed and valid
    pow := core.NewProofOfWork(input.Bits, []byte(info.Nonce), input.Proof)
    if ! pow.Validate() {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotAcceptable)
        _ = json.NewEncoder(w).Encode(StatusError("proof-of-work cannot be validated"))
        return
    }

    // We can generate an accept path
    path, err := is.GenerateAcceptResponsePath(info.Addr, info.Checksum)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    // Allow 30 minutes for incoming body message
    to := time.Now()
    to.Add(time.Minute * 30)

    ret := OutputHeaderType{
        Error: false,
        Status: BODY_ACCEPT,
        Description: "Accepting body for this header",
        BodyAccept: &BodyAcceptType{
            Path: "/incoming/" + path,
            Timeout: to.Format(time.RFC3339),
        },
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(ret)
}


func handleAccept(w http.ResponseWriter, req *http.Request, info *incoming.IncomingInfoType) {

}




package handler

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/core/container"
    "github.com/jaytaph/mailv2/core/message"
    "github.com/jaytaph/mailv2/server/incoming"
    "github.com/jaytaph/mailv2/core/utils"
    "net/http"
    "time"
)


func PostMessageBody(w http.ResponseWriter, req *http.Request) {
    is := container.GetIncomingService()
    path := mux.Vars(req)["id"]

    // @TODO: We need to check if path is actually an UUID, otherwise we could browse through the whole redis DB
    info, err := is.GetIncomingInfo(path)
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

    switch info.Type {
    case incoming.ACCEPT:
        handleAccept(w, req, info)
        return
    case incoming.PROOF_OF_WORK:
        handlePow(w, req, info)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    _ = json.NewEncoder(w).Encode(StatusError("unknown incoming type for this request"))
}

func handlePow(w http.ResponseWriter, req *http.Request, info *incoming.IncomingInfoType) {
    is := container.GetIncomingService()
    decoder := json.NewDecoder(req.Body)

    // Decode JSON
    var input message.MessageProofOfWork
    err := decoder.Decode(&input)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    if ! utils.ValidateProofOfWork(info.Bits, []byte(info.Nonce), input.Data) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotAcceptable)
        _ = json.NewEncoder(w).Encode(StatusError("proof-of-work cannot be validated"))
        return
    }

    // POW has been done. We can create an accept path
    path, err := is.GenerateAcceptPath(info.Email, info.Checksum)
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




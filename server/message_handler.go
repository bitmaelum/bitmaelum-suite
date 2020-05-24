package server

import (
    "encoding/json"
    "github.com/jaytaph/mailv2/container"
    http_status "github.com/jaytaph/mailv2/http"
    "github.com/jaytaph/mailv2/message"
    "math/rand"
    "net/http"
    "time"
)

const (
    PROOF_OF_WORK   string = "proof_of_work"
    BODY_ACCEPT     string = "body_accept"
)

type ProofOfWorkType struct {
    Bits        int         `json:"bits"`
    Nonce       string      `json:"nonce"`
    Path        string      `json:"path"`
    Timeout     string      `json:"timeout"`
}

type BodyAcceptType struct {
    Path        string      `json:"path"`
    Timeout     string      `json:"timeout"`
}

type OutputHeaderType struct {
    Error           bool                `json:"error"`
    Status          string              `json:"status"`
    Description     string              `json:"description"`
    ProofOfWork     *ProofOfWorkType    `json:"proof_of_work,omitempty"`
    BodyAccept      *BodyAcceptType     `json:"body_accept,omitempty"`
}

func PostMessageHeader(w http.ResponseWriter, req *http.Request) {
    is := container.GetIncomingService()

    decoder := json.NewDecoder(req.Body)

    // Decode JSON
    var input message.MessageHeader
    err := decoder.Decode(&input)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(http_status.StatusError("Malformed JSON: " + err.Error()))
        return
    }

    // Validate incoming header
    err = message.ValidateHeader(input)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(http_status.StatusError(err.Error()))
        return
    }

    // Check if we need proof of work.
    if needsProofOfWork(input) {
        // Generate proof-of-work data
        path, nonce, err := is.GeneratePowPath(input.From.Email, 24)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            _ = json.NewEncoder(w).Encode(http_status.StatusError(err.Error()))
            return
        }

        // Allow 30 minutes for proof-of-work
        to := time.Now()
        to.Add(time.Minute * 30)

        pow := &ProofOfWorkType{
            Bits: 24,
            Nonce: nonce,
            Path: "/incoming/" + path,
            Timeout: to.Format(time.RFC3339),
        }

        ret := OutputHeaderType{
            Error: false,
            Status: PROOF_OF_WORK,
            Description: "A proof of work is needed before we will accept this message",
            ProofOfWork: pow,
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _ = json.NewEncoder(w).Encode(ret)
        return
    }

    // No proof-of-work, generate accept path
    path, err := is.GenerateAcceptPath(input.From.Email)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(http_status.StatusError(err.Error()))
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

func needsProofOfWork(header message.MessageHeader) bool {
    // @TODO: We probably want to use different metrics to check if we need to do proof-of-work
    return rand.Intn(10) < 5
}

func PostMessageBody(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode([]byte{})
}




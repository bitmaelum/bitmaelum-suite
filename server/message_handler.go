package server

import (
    "encoding/json"
    "github.com/google/uuid"
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
    Load    int      `json:"load"`
    Data    []byte   `json:"data"`
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
    decoder := json.NewDecoder(req.Body)

    var input message.MessageHeader
    err := decoder.Decode(&input)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(http_status.StatusError("Malformed JSON: " + err.Error()))
        return
    }

    err = message.ValidateHeader(input)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(http_status.StatusError(err.Error()))
        return
    }

    pow, err := needsProofOfWork(input)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(http_status.StatusError(err.Error()))
        return
    }

    // Seems like this message needs to do some proof of work before continuing
    if pow != nil {
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

    path, err := uuid.NewRandom()
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
            Path: "/incoming/" + path.String(),
            Timeout: to.Format(time.RFC3339),
        },
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(ret)

}

func needsProofOfWork(header message.MessageHeader) (*ProofOfWorkType, error) {
    // @TODO: We probably want to use different metrics to check if we need to do proof-of-work
    if rand.Intn(10) < 5 {
        rnd := make([]byte, 20)
        _, err := rand.Read(rnd)
        if err != nil {
            return nil, err
        }

        pow := &ProofOfWorkType{
            Load: 10,
            Data: rnd,
        }
        return pow, nil
    }

    return nil, nil;
}

func PostMessageBody(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode([]byte{})
}




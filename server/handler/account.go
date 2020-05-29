package handler

import (
    "encoding/json"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/core/container"
    "github.com/jaytaph/mailv2/core/utils"
    "net/http"
)

type InputCreateAccount struct {
    Mailbox     string  `json:"mailbox"`
    PublicKey   string  `json:"public_key"`
    ProofOfWork struct {
        Bits     int     `json:"bits"`
        Proof    uint64  `json:"proof"`
    } `json:"proof_of_work"`
}

func CreateAccount(w http.ResponseWriter, req *http.Request) {
    if ! config.Server.Account.Registration {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusForbidden)
        _ = json.NewEncoder(w).Encode(StatusError("public registration not available"))
        return
    }

    var input InputCreateAccount
    err := DecodeBody(w, req.Body, &input)
    if err != nil {
        return
    }

    // Check proof of work first
    if input.ProofOfWork.Bits < config.Server.Account.ProofOfWork {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusErrorf("Proof of work must be at least %d bits", config.Server.Account.ProofOfWork))
        return
    }
    if ! utils.ValidateProofOfWork(input.ProofOfWork.Bits, []byte(input.Mailbox), input.ProofOfWork.Proof) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusError("Proof of work incorrect"))
        return
    }

    as := container.GetAccountService()
    err = as.CreateAccount(input.Mailbox, input.PublicKey)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(StatusOk("mailbox created"))
}

func RetrieveAccount(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(StatusOk("this is your account"))
    return
}

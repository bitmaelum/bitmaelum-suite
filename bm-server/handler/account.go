package handler

import (
    "encoding/json"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/bitmaelum/bitmaelum-server/core/container"
    "net/http"
)

type InputCreateAccount struct {
    Addr        core.HashAddress    `json:"address"`
    Token       string              `json:"token"`
    PublicKey   string              `json:"public_key"`
    ProofOfWork struct {
        Bits     int                `json:"bits"`
        Proof    uint64             `json:"proof"`
    } `json:"proof_of_work"`
}

// Create account handler
func CreateAccount(w http.ResponseWriter, req *http.Request) {

    // Only allow registration when enabled in the configuration
    if ! config.Server.Accounts.Registration {
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
    if input.ProofOfWork.Bits < config.Server.Accounts.ProofOfWork {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusErrorf("Proof of work must be at least %d bits", config.Server.Accounts.ProofOfWork))
        return
    }
    pow := core.NewProofOfWork(input.ProofOfWork.Bits, []byte(input.Addr), input.ProofOfWork.Proof)
    if ! pow.Validate() {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusError("Proof of work incorrect"))
        return
    }

    // @TODO: check if account already exists?

    // Create account
    as := container.GetAccountService()
    err = as.CreateAccount(input.Addr, input.PublicKey)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(StatusOk("mailbox created"))
}

// Retrieve account information
func RetrieveAccount(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(StatusOk("this is your account"))
    return
}

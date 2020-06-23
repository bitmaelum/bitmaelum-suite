package handler

import (
    "crypto/subtle"
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

    // Check if token exists for the given address
    is := container.GetInviteService()
    registeredToken, err := is.GetInvite(input.Addr)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusError("No token found for this address."))
        return
    }
    if subtle.ConstantTimeCompare([]byte(registeredToken), []byte(input.Token)) != 1 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusError("Incorrect token found for this address."))
        return
    }

    // Check if account exists
    as := container.GetAccountService()
    if as.AccountExists(input.Addr) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(StatusError("Account already exists."))
        return
    }

    // All clear. Create account
    err = as.CreateAccount(input.Addr, input.PublicKey)
    if err != nil {
        sendBadRequest(w, err)
        return
    }

    // Done with the invite, let's remove
    _ = is.RemoveInvite(input.Addr)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(StatusOk("BitMaelum account has been successfully created."))
}

// Retrieve account information
func RetrieveAccount(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(StatusOk("this is your account"))
    return
}

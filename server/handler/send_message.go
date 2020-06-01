package handler

import (
    "bytes"
    "crypto/sha256"
    "encoding/json"
    "github.com/jaytaph/mailv2/core/container"
    "github.com/jaytaph/mailv2/core/message"
    "io/ioutil"
    "math/rand"
    "net/http"
    "time"
)

// Handler when a message header is posted
func SendMessage(w http.ResponseWriter, req *http.Request) {

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

package core

import (
    "crypto/sha256"
    "encoding/hex"
)

func HashEmail(email string) string {
    sum := sha256.Sum256([]byte(email))
    return hex.EncodeToString(sum[:])
}


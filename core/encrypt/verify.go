package encrypt

import (
    "crypto/subtle"
)

// Verify if hash compares against the signature of the message
func Verify(key interface{}, message []byte, hash []byte) (bool, error) {
    computedHash, err := Sign(key, message)
    if err != nil {
        return false, err
    }

    return subtle.ConstantTimeCompare(computedHash, hash) == 1, nil
}

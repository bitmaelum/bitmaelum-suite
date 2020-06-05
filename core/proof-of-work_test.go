package core

import (
    "testing"
)

func Test_ProofOfWork(t *testing.T) {
    if ProofOfWork(8, []byte("john@example.org")) != 456 {
        t.Error()
    }

    if ProofOfWork(16, []byte("john@example.org")) != 28544 {
        t.Error()
    }
}

func Test_ValidateProofOfWork(t *testing.T) {
    if ! ValidateProofOfWork(8, []byte("john@example.org"), 456) {
        t.Error()
    }

    if ! ValidateProofOfWork(16, []byte("john@example.org"), 28544) {
        t.Error()
    }
}


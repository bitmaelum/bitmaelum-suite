package message

import (
    "errors"
    "github.com/jaytaph/mailv2/core/utils"
)

func ValidateHeader(input MessageHeader) error {

    if input.From.ProofOfWork.Bits < 20 {
        return errors.New("proof of work is too low")
    }

    if ! utils.ValidateProofOfWork(int(input.From.ProofOfWork.Bits), []byte(input.From.Email), input.From.ProofOfWork.Nonce) {
        return errors.New("invalid proof of work")
    }

    // Check if there is a mailbox for this account.

    // Other checks that we can / must do

    return nil
}

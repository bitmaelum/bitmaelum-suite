package examples

import (
    "fmt"
    "github.com/jaytaph/mailv2/utils"
)

func main() {
    var data []byte = []byte("joshua@noxlogic.nl")

    nonce := utils.ProofOfWork(20, data)
    fmt.Printf("\n\nNonce found: %#v\n", nonce)

    if utils.ValidateProofOfWork(20, data, nonce) {
        fmt.Printf("Valid POW")
    } else {
        fmt.Printf("Invalid POW")
    }
}

package main

import (
    "fmt"
    "github.com/jaytaph/mailv2/utils"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Printf("Usage: %s <bits> <nonce>", os.Args[0])
        os.Exit(1)
    }

    bits, _ := strconv.Atoi(os.Args[1])
    nonce := os.Args[2]

    fmt.Printf("Working on %d bits proof...\n", bits)
    pow := utils.ProofOfWork(bits, []byte(nonce))

    fmt.Printf("Proof: %d\n", pow)
}

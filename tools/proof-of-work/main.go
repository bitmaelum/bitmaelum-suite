package main

import (
    "fmt"
    "github.com/jaytaph/mailv2/core/utils"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Printf("Usage: %s <bits> <data>", os.Args[0])
        os.Exit(1)
    }

    bits, _ := strconv.Atoi(os.Args[1])
    data := os.Args[2]

    fmt.Printf("Working on %d bits proof...\n", bits)
    pow := utils.ProofOfWork(bits, []byte(data))

    fmt.Printf("Proof: %d\n", pow)
}

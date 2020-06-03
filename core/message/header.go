package message

import (
    "github.com/jaytaph/mailv2/core"
)

type Header struct {
    From struct {
        Addr        core.HashAddress    `json:"address"`
        PublicKey   string              `json:"public_key"`
        ProofOfWork struct {
            Bits    int     `json:"bits"`
            Proof   uint64  `json:"proof"`
        } `json:"proof_of_work"`
    } `json:"from"`
    To struct {
        Addr    core.HashAddress    `json:"address"`
    } `json:"to"`
    Catalog struct {
        Size        uint64          `json:"size"`
        Checksum    []ChecksumType  `json:"checksum"`
        Crypto      string          `json:"crypto"`
        Key         []byte          `json:"key"`
        Iv          []byte          `json:"iv"`
    } `json:"catalog"`
}

type ProofOfWorkType struct {
   Bits     int     `json:"bits"`
   Proof    uint64  `json:"proof"`
}

type ChecksumType struct {
   Hash    string  `json:"hash"`
   Value   string  `json:"value"`
}

package message

import (
    "github.com/jaytaph/mailv2/core"
)

type Message struct {
    Header  Header
    Catalog Catalog
}

func LoadMessage(addr core.Address, id string) {
}

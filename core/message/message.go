package message

import (
    "github.com/bitmaelum/bitmaelum-server/core"
)

type Message struct {
    Header  Header
    Catalog Catalog
}

func LoadMessage(addr core.Address, id string) {
}

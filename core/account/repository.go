package account

import "github.com/jaytaph/mailv2/core"

type Repository interface {
    // Account management
    Create(addr core.HashAddress) error
    Exists(addr core.HashAddress) bool

    StorePubKey(addr core.HashAddress, data []byte) error
    FetchPubKey(addr core.HashAddress) ([]byte, error)

    CreateBox(addr core.HashAddress, box, description string, quota int) error
    ExistsBox(addr core.HashAddress, box string) bool
    DeleteBox(addr core.HashAddress, box string) error
}

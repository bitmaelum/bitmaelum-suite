package account

import (
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/message"
)

type Repository interface {
    // Account management
    Create(addr core.HashAddress) error
    Exists(addr core.HashAddress) bool

    // Public key
    StorePubKey(addr core.HashAddress, key string) error
    FetchPubKeys(addr core.HashAddress) ([]string, error)

    // Box related functions
    CreateBox(addr core.HashAddress, box, description string, quota int) error
    ExistsBox(addr core.HashAddress, box string) bool
    DeleteBox(addr core.HashAddress, box string) error
    GetBox(addr core.HashAddress, box string) (*message.MailBoxInfo, error)
    FindBox(addr core.HashAddress, query string) ([]message.MailBoxInfo, error)

    // Message boxes
    FindMessages(addr core.HashAddress, box string, offset, limit int) ([]message.MessageInfo, error)
}

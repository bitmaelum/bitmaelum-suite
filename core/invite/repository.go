package invite

import (
    "github.com/jaytaph/mailv2/core"
    "time"
)

type Repository interface {
    CreateInvite(addr core.HashAddress, expiry time.Duration) (string, error)
    GetInvite(addr core.HashAddress) (string, error)
    RemoveInvite(addr core.HashAddress) error
}

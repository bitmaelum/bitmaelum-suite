package resolve

import "github.com/jaytaph/mailv2/core"

type Repository interface {
    Resolve(addr core.HashAddress) (*ResolveInfo, error)
    Upload(addr core.HashAddress, pubKey, resolveAddress, signature string) error
}

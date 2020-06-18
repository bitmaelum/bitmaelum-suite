package resolve

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/account/server"
	"github.com/sirupsen/logrus"
)

type localRepo struct {
    as *server.Service
}

func NewLocalRepository(s *server.Service) Repository {
    return &localRepo{
        as: s,
    }
}

func (r *localRepo) Resolve(addr core.HashAddress) (*ResolveInfo, error) {
    logrus.Trace("local repository cache is not available.")

    return nil, errors.New("key not found in local cache")
}

func (r *localRepo) Upload(addr core.HashAddress, pubKey, address, signature string) error {
    return nil
}

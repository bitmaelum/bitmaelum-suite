package keys

import (
    "errors"
    "github.com/sirupsen/logrus"
)

type remoteRepo struct {
}

func NewRemoteRepository() Repository {
    return &remoteRepo{}
}

func (r *remoteRepo) Retrieve(hash string) (string, error) {
    logrus.Warning("Remote repositories are not implemented yet")

    return "", errors.New("Key not found")
}

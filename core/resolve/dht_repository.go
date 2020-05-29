package resolve

import (
    "errors"
    "github.com/sirupsen/logrus"
)

type dhtRepo struct {
}

func NewDHTRepository() Repository {
    return &dhtRepo{}
}

func (r *dhtRepo) Retrieve(hash string) (*ResolveInfo, error) {
    logrus.Trace("DHT is not implemented yet")

    return nil, errors.New("key not found in DHT")
}

func (r *dhtRepo) Upload(hash, pubKey, address, signature string) error {
    return nil
}

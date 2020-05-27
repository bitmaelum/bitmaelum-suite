package keys

import (
    "errors"
    "github.com/sirupsen/logrus"
)

type dhtRepo struct {
}

func NewDHTRepository() Repository {
    return &dhtRepo{}
}

func (r *dhtRepo) Retrieve(hash string) (string, error) {
    logrus.Warning("DHT is not implemented yet")

    return "", errors.New("Key not found")
}

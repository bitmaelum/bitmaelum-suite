package incoming

import "time"

type Reader interface {
	Has(path string) (bool, error)
	Get(path string) ([]byte, error)
}

type Writer interface {
	Create(path string, value []byte, expiry time.Duration) error
	Remove(path string) error
}

type Repository interface {
	Reader
	Writer
}

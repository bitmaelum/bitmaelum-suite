package incoming

import "time"

// Reader is the interface to read from the incoming repository
type Reader interface {
	Has(path string) (bool, error)
	Get(path string) ([]byte, error)
}

// Writer is the interface to write to the incoming repository
type Writer interface {
	Create(path string, value []byte, expiry time.Duration) error
	Remove(path string) error
}

// Repository is the complete interface for reading/writing
type Repository interface {
	Reader
	Writer
}

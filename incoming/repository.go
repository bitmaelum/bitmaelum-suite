package incoming

import "time"

type Reader interface {
    Has(path string) (bool, error)
    Get(path string) (string, error)
}

type Writer interface {
    Set(path string, value string, expiry time.Duration) error
}

type IncomingRepository interface {
    Reader
    Writer
}


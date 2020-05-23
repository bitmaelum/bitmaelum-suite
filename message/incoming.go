package message

type Reader interface {
    Has(path string) (bool, error)
    Get(path string) (string, error)
}

type Writer interface {
    Set(path string, value string)
}

type Repository interface {
    Reader
    Writer
}

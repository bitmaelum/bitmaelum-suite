package keys

type Repository interface {
    Retrieve(hash string) (string, error)
}

package account

type Repository interface {
    // Account management
    Create(hash string) error
    Exists(hash string) bool

    StorePubKey(hash string, path string, data []byte) error
    FetchPubKey(hash string, path string) ([]byte, error)

    CreateBox(hash string, box string, description string, quota int) error
    ExistsBox(hash string, box string) bool
    DeleteBox(hash string, box string) error
}


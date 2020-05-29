package resolve

type Repository interface {
    Retrieve(hash string) (*ResolveInfo, error)
    Upload(hash, pubKey, address, signature string) error
}

package resolve

type Repository interface {
    Retrieve(hash string) (*ResolveInfo, error)
    Upload(hash, pubKey, resolveAddress, signature string) error
}

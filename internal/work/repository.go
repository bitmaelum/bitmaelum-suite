package work

type Repository interface {
	GetName() string
	GetWorkOutput() map[string]interface{}
	GetWorkProofOutput() map[string]interface{}
	ValidateWork(data []byte) bool
	Work()
}

// GetPreferredWork will fetch work by taken the preference into account
func GetPreferredWork(preferences []string) (Repository, error){
	// The preference list is just that: the preference of the client. The server ultimately decides what kind of
	// work needs to be done.

	// Currently, only POW is supported
	return NewPow()
}

package incoming

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/google/uuid"
	"time"
)

const (
	// ProofOfWork constant
	ProofOfWork string = "pow"
	// Accept constant
	Accept string = "accept"
)

// Service representing an incoming service
type Service struct {
	repo Repository
}

// InfoType is a structure
type InfoType struct {
	Type     string              `json:"type"`
	Addr     address.HashAddress `json:"address"`
	Nonce    string              `json:"nonce,omitempty"`
	Bits     int                 `json:"bits,omitempty"`
	Checksum []byte              `json:"checksum"`
}

// NewIncomingService Create new incoming service
func NewIncomingService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GenerateAcceptResponsePath generates an accept response path
func (is *Service) GenerateAcceptResponsePath(addr address.HashAddress, checksum []byte) (string, error) {
	data := &InfoType{
		Type:     Accept,
		Addr:     addr,
		Checksum: checksum,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	path, err := is.generatePath()
	if err != nil {
		return "", err
	}

	_ = is.repo.Create(path, jsonData, time.Duration(time.Minute*30))

	return path, nil
}

// GeneratePowResponsePath generates an proof-of-work response
func (is *Service) GeneratePowResponsePath(addr address.HashAddress, bits int, checksum []byte) (string, string, error) {
	rnd := make([]byte, 32)
	_, err := rand.Read(rnd)
	if err != nil {
		return "", "", err
	}
	nonce := base64.StdEncoding.EncodeToString(rnd)

	data := &InfoType{
		Type:     ProofOfWork,
		Addr:     addr,
		Nonce:    nonce,
		Bits:     bits,
		Checksum: checksum,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	path, err := is.generatePath()
	if err != nil {
		return "", "", err
	}

	_ = is.repo.Create(path, jsonData, time.Duration(time.Minute*30))

	return path, nonce, nil
}

// RemovePath removes an incoming path
func (is *Service) RemovePath(path string) error {
	_ = is.repo.Remove(path)

	return nil
}

// Generates a new incoming path for either proof-of-work or accept responses
func (is *Service) generatePath() (string, error) {
	p, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return p.String(), nil
}

// GetIncomingPath gets incoming info
func (is *Service) GetIncomingPath(path string) (*InfoType, error) {
	found, err := is.repo.Has(path)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}

	data, err := is.repo.Get(path)
	if err != nil {
		return nil, err
	}

	incomingInfo := InfoType{}
	err = json.Unmarshal(data, &incomingInfo)
	if err != nil {
		return nil, err
	}

	return &incomingInfo, nil
}

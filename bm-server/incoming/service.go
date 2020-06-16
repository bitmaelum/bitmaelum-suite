package incoming

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/bitmaelum/bitmaelum-server/core"
    "time"
)

const (
    PROOF_OF_WORK   string = "pow"
    ACCEPT          string = "accept"
)

type Service struct {
    repo Repository
}

type IncomingInfoType struct {
    Type        string              `json:"type"`
    Addr        core.HashAddress    `json:"address"`
    Nonce       string              `json:"nonce,omitempty"`
    Bits        int                 `json:"bits,omitempty"`
    Checksum    []byte              `json:"checksum"`
}

// Create new incoming service
func NewIncomingService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

// Generate an accept response path
func (is *Service) GenerateAcceptResponsePath(addr core.HashAddress, checksum []byte) (string, error) {
    data := &IncomingInfoType{
        Type: ACCEPT,
        Addr: addr,
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

    _ = is.repo.Create(path, jsonData, time.Duration(time.Minute * 30))

    return path, nil
}

// Generate an proof-of-work response
func (is *Service) GeneratePowResponsePath(addr core.HashAddress, bits int, checksum []byte) (string, string, error) {
    rnd := make([]byte, 32)
    _, err := rand.Read(rnd)
    if err != nil {
        return "", "", err
    }
    nonce := base64.StdEncoding.EncodeToString(rnd)

    data := &IncomingInfoType{
        Type: PROOF_OF_WORK,
        Addr: addr,
        Nonce: nonce,
        Bits: bits,
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

    _ = is.repo.Create(path, jsonData, time.Duration(time.Minute * 30))

    return path, nonce, nil
}

// Remove an incoming path
func (is *Service) RemovePath(path string) error {
    _ = is.repo.Remove(path)

    return nil
}

// Generates a new incoming path for either proof-of-work or accept responses
func (is *Service) generatePath() (string, error) {
    path, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }

    return path.String(), nil
}

// Get incoming info
func (is *Service) GetIncomingPath(path string) (*IncomingInfoType, error) {
    found, err := is.repo.Has(path)
    if err != nil {
        return nil, err
    } else if ! found {
        return nil, nil
    }

    data, err := is.repo.Get(path)
    if err != nil {
        return nil, err
    }

    incomingInfo := IncomingInfoType{}
    err = json.Unmarshal(data, &incomingInfo)
    if err != nil {
        return nil, err
    }

    return &incomingInfo, nil
}

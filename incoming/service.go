package incoming

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "github.com/google/uuid"
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
    Type        string      `json:"type"`
    Email       string      `json:"email"`
    Nonce       string      `json:"nonce,omitempty"`
    Bits        int         `json:"bits,omitempty"`
    Checksum    []byte      `json:"checksum"`
}

func NewIncomingService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

func (is *Service) GenerateAcceptPath(email string, checksum []byte) (string, error) {
    data := &IncomingInfoType{
        Type: ACCEPT,
        Email: email,
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

func (is *Service) GeneratePowPath(email string, bits int, checksum []byte) (string, string, error) {
    rnd := make([]byte, 32)
    _, err := rand.Read(rnd)
    if err != nil {
        return "", "", err
    }
    nonce := base64.StdEncoding.EncodeToString(rnd)

    data := &IncomingInfoType{
        Type: PROOF_OF_WORK,
        Email: email,
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

func (is *Service) RemovePath(path string) error {
    _ = is.repo.Remove(path)

    return nil
}

func (is *Service) generatePath() (string, error) {
    path, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }

    return path.String(), nil
}

func (is *Service) GetIncomingInfo(path string) (*IncomingInfoType, error) {
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

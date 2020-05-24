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

type redisType struct {
    Type    string  `json:"type"`
    Email   string  `json:"email"`
    Nonce   string  `json:"nonce,omitempty"`
    Bits    int     `json:"bits,omitempty"`
}

func NewIncomingService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

func (is *Service) GenerateAcceptPath(email string) (string, error) {
    data := &redisType{
        Type: ACCEPT,
        Email: email,
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

func (is *Service) GeneratePowPath(email string, bits int) (string, string, error) {
    rnd := make([]byte, 32)
    _, err := rand.Read(rnd)
    if err != nil {
        return "", "", err
    }
    nonce := base64.StdEncoding.EncodeToString(rnd)

    data := &redisType{
        Type: PROOF_OF_WORK,
        Email: email,
        Nonce: nonce,
        Bits: bits,
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

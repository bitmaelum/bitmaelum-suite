package resolve

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/bitmaelum/bitmaelum-suite/core/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"net/url"
	"strings"
)

// Service represents a resolver service tied to a specific repository
type Service struct {
	repo Repository
}

// Info is a structure returned by the external resolver system
type Info struct {
	Hash      string `json:"hash"`
	PublicKey string `json:"public_key"`
	Address   string `json:"address"`
}

// KeyRetrievalService initialises a key retrieval service.
func KeyRetrievalService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// IsLocal returns true when the resolve address of the info is our own server address
func (info *Info) IsLocal() bool {
	localHost, localPort, err := getHostPort(config.Server.Server.Name)
	if err != nil {
		return false
	}
	infoHost, infoPort, err := getHostPort(info.Address)
	if err != nil {
		return false
	}

	return localHost == infoHost && localPort == infoPort
}

func getHostPort(hostport string) (string, string, error) {
	// We need schema otherwise url.Parse does not work
	if !strings.HasPrefix(hostport, "http://") && !strings.HasPrefix(hostport, "https://") {
		hostport = "https://" + hostport
	}

	info, err := url.Parse(hostport)
	if err != nil {
		return "", "", err
	}

	if info.Port() == "" {
		info.Host = info.Host + ":2424"
	}

	return info.Hostname(), info.Port(), nil
}

// Resolve resolves an address.
func (s *Service) Resolve(addr address.HashAddress) (*Info, error) {
	return s.repo.Resolve(addr)
}

// UploadInfo uploads resolve information to a service.
func (s *Service) UploadInfo(info account.Info, resolveAddress string) error {

	// @TODO: We maybe should sign with a different algo? Otherwise we use the same one for all systems
	privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
	if err != nil {
		return err
	}

	// Sign resolve address
	hash := sha256.Sum256([]byte(resolveAddress))
	signature, err := encrypt.Sign(privKey, hash[:])
	if err != nil {
		return err
	}

	h, err := address.NewHash(info.Address)
	if err != nil {
		return err
	}

	// And upload
	return s.repo.Upload(
		*h,
		string(info.PubKey),
		resolveAddress,
		hex.EncodeToString(signature),
	)
}

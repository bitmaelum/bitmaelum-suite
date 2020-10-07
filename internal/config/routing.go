package config

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/afero"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/hkdf"
)

// Routing holds routing configuration for the mail server
type Routing struct {
	RoutingID  string           `json:"routing_id"`
	PrivateKey bmcrypto.PrivKey `json:"private_key"`
	PublicKey  bmcrypto.PubKey  `json:"public_key"`
}

// ReadRouting will read the routing file and merge it into the server configuration
func ReadRouting(p string) error {
	f, err := fs.Open(p)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	Server.Routing = &Routing{}
	err = json.Unmarshal(data, Server.Routing)
	if err != nil {
		return err
	}

	return nil
}

// SaveRouting will save the routing into a file. It will overwrite if exists
func SaveRouting(p string, routing *Routing) error {
	data, err := json.MarshalIndent(routing, "", "  ")
	if err != nil {
		return err
	}

	err = fs.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, p, data, 0600)
}

// Generate generates a new routing structure
func GenerateRouting() (string, *Routing, error) {
	// Generate large enough random string
	e, err := bip39.NewEntropy(192)
	if err != nil {
		return "", nil, err
	}

	// Generate seed words
	seed, err := bip39.NewMnemonic(e)
	if err != nil {
		return "", nil, err
	}

	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, e, []byte{}, []byte{})
	expbuf := make([]byte, 32)
	_, err = io.ReadFull(rd, expbuf)
	if err != nil {
		return "", nil, err
	}

	// Generate keypair
	r := ed25519.NewKeyFromSeed(expbuf[:32])
	privKey, err := bmcrypto.NewPrivKeyFromInterface(r)
	if err != nil {
		return "", nil, err
	}
	pubKey, err := bmcrypto.NewPubKeyFromInterface(r.Public())
	if err != nil {
		return "", nil, err
	}

	return seed, &Routing{
		RoutingID:  hash.New(hex.EncodeToString(expbuf)).String(),
		PrivateKey: *privKey,
		PublicKey:  *pubKey,
	}, nil
}

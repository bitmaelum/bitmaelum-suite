package routing

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
)

type RoutingType struct {
	RoutingID  string           `json:"routing_id"`
	PrivateKey bmcrypto.PrivKey `json:"private_key"`
	PublicKey  bmcrypto.PubKey  `json:"public_key"`
}


// ReadRouting will read the routing file
func ReadRouting(p string) (*RoutingType, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	routing := &RoutingType{}
	err = json.Unmarshal(data, routing)
	if err != nil {
		return nil, err
	}

	return routing, nil
}

// SaveRouting will save the routing into a file. It will overwrite if exists
func SaveRouting(p string, routing *RoutingType) error {
	data, err := json.MarshalIndent(routing, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(p), 755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, data, 0600)
}

// Generate generates a new routing structure
func Generate() (*RoutingType, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256([]byte(id.String()))
	routingID := hex.EncodeToString(sum[:])


	privKey, pubKey, err := encrypt.GenerateKeyPair(bmcrypto.KeyTypeRSA)
	if err != nil {
		return nil, err
	}

	return &RoutingType{
		RoutingID:  routingID,
		PrivateKey: *privKey,
		PublicKey:  *pubKey,
	}, nil
}

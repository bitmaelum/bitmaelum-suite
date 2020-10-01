package internal

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// AccountInfo represents client account information
type AccountInfo struct {
	Default bool   `json:"default"` // Is this the default account
	Address string `json:"address"` // The address of the account

	Name     string            `json:"name"`     // Full name of the user
	Settings map[string]string `json:"settings"` // Additional settings that can be user-defined

	// Communication and encryption information
	PrivKey   bmcrypto.PrivKey        `json:"priv_key"`        // PEM encoded private key
	PubKey    bmcrypto.PubKey         `json:"pub_key"`         // PEM encoded public key
	Pow       proofofwork.ProofOfWork `json:"proof,omitempty"` // Proof of work
	RoutingID string                  `json:"routing_id"`      // ID of the routing used
}

// OrganisationInfo represents a organisation configuration for a server
type OrganisationInfo struct {
	Addr        string                        `json:"addr"`          // org part from the bitmaelum address
	Name        string                        `json:"name"`          // Full name of the organisation
	PrivKey     bmcrypto.PrivKey              `json:"priv_key"`      // PEM encoded private key
	PubKey      bmcrypto.PubKey               `json:"pub_key"`       // PEM encoded public key
	Pow         proofofwork.ProofOfWork       `json:"pow,omitempty"` // Proof of work
	Validations []organisation.ValidationType `json:"validations"`   // Validations
}

// RoutingInfo represents a routing configuration for a server
type RoutingInfo struct {
	RoutingID string                  `json:"routing_id"`    // ID
	PrivKey   bmcrypto.PrivKey        `json:"priv_key"`      // PEM encoded private key
	PubKey    bmcrypto.PubKey         `json:"pub_key"`       // PEM encoded public key
	Pow       proofofwork.ProofOfWork `json:"pow,omitempty"` // Proof of work
	Route     string                  `json:"route"`         // Route to server
}

func InfoToOrg(info OrganisationInfo) (*organisation.Organisation, error) {
	a, err := address.NewOrgHash(info.Addr)
	if err != nil {
		return nil, err
	}

	return &organisation.Organisation{
		Addr:       *a,
		Name:       info.Name,
		PublicKey:  info.PubKey,
		Validation: info.Validations,
	}, nil
}

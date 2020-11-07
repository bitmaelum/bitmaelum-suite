package internal

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// RoutingInfo represents a routing configuration for a server
type RoutingInfo struct {
	RoutingID string                   `json:"routing_id"`    // ID
	PrivKey   bmcrypto.PrivKey         `json:"priv_key"`      // PEM encoded private key
	PubKey    bmcrypto.PubKey          `json:"pub_key"`       // PEM encoded public key
	Pow       *proofofwork.ProofOfWork `json:"pow,omitempty"` // Proof of work
	Route     string                   `json:"route"`         // Route to server
}

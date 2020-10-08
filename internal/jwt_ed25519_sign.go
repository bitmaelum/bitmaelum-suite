package internal

import (
	"crypto/ed25519"

	"github.com/vtolstov/jwt-go"
)

// taken from https://blainsmith.com/articles/signing-jwts-with-gos-crypto-ed25519/

// SigningMethodEdDSA ed25519 signing method for JWT token signing
type SigningMethodEdDSA struct{}

// Alg returns the algorithm used
func (m *SigningMethodEdDSA) Alg() string {
	return "EdDSA"
}

// Verify will verify a signature
func (m *SigningMethodEdDSA) Verify(signingString string, signature string, key interface{}) error {
	var err error

	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	var ed25519Key ed25519.PublicKey
	var ok bool
	if ed25519Key, ok = key.(ed25519.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	if ok := ed25519.Verify(ed25519Key, []byte(signingString), sig); !ok {
		return jwt.ErrECDSAVerification
	}

	return nil
}

// Sign wil sign a string
func (m *SigningMethodEdDSA) Sign(signingString string, key interface{}) (str string, err error) {
	var ed25519Key ed25519.PrivateKey
	var ok bool
	if ed25519Key, ok = key.(ed25519.PrivateKey); !ok {
		return "", jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PrivateKeySize {
		return "", jwt.ErrInvalidKey
	}

	// Sign the string and return the encoded result
	sig := ed25519.Sign(ed25519Key, []byte(signingString))
	return jwt.EncodeSegment(sig), nil
}

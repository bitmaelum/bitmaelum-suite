package bmcrypto

import (
	"crypto/ed25519"
	"crypto/sha256"
	"io"

	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/hkdf"
)

// GenerateKeypairFromSeed generates a keypair based on the given seed
func GenerateKeypairFromSeed(seed string) (*PrivKey, *PubKey, error) {
	e, err := bip39.MnemonicToByteArray(seed, true)
	if err != nil {
		return nil, nil, err
	}

	return genKey(e)
}

// GenerateKeypairWithSeed generates a seed, and a keypair that can be generated through the same seed again.
func GenerateKeypairWithSeed() (string, *PrivKey, *PubKey, error) {
	// Generate large enough random string
	e, err := bip39.NewEntropy(192)
	if err != nil {
		return "", nil, nil, err
	}

	// Generate seed words
	seed, err := bip39.NewMnemonic(e)
	if err != nil {
		return "", nil, nil, err
	}

	privKey, pubKey, err := genKey(e)
	if err != nil {
		return "", nil, nil, err
	}

	return seed, privKey, pubKey, nil
}

func genKey(e []byte) (*PrivKey, *PubKey, error) {
	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, e, []byte{}, []byte{})
	expbuf := make([]byte, 32)
	_, err := io.ReadFull(rd, expbuf)
	if err != nil {
		return nil, nil, err
	}

	// Generate keypair
	r := ed25519.NewKeyFromSeed(expbuf[:32])
	privKey, err := NewPrivKeyFromInterface(r)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := NewPubKeyFromInterface(r.Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

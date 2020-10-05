package bmcrypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha512"
	"errors"
	"math/big"

	"golang.org/x/crypto/curve25519"
)

// KeyExchange exchange a message given the Private and other's Public Key
func KeyExchange(privK PrivKey, pubK PubKey) ([]byte, error) {
	switch privK.Type {
	case KeyTypeECDSA:
		k, _ := pubK.K.(*ecdsa.PublicKey).Curve.ScalarMult(pubK.K.(*ecdsa.PublicKey).X, pubK.K.(*ecdsa.PublicKey).Y, privK.K.(*ecdsa.PrivateKey).D.Bytes())

		return k.Bytes(), nil

	case KeyTypeED25519:
		x25519priv := EdPrivToX25519(privK.K.(ed25519.PrivateKey))
		x25519pub := EdPubToX25519(pubK.K.(ed25519.PublicKey))
		return curve25519.X25519(x25519priv, x25519pub)
	}

	return nil, errors.New("unknown key type for key exchange")
}

//EdPrivToX25519 converts a ed25519 PrivateKey to a X25519 Private Key
func EdPrivToX25519(privateKey ed25519.PrivateKey) []byte {
	h := sha512.New()
	_, _ = h.Write(privateKey[:32])
	digest := h.Sum(nil)
	h.Reset()

	// From https://cr.yp.to/ecdh.html (I don't think this is really needed in this case)
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	return digest[:32]
}

//From https://github.com/FiloSottile/age/blob/bbab440e198a4d67ba78591176c7853e62d29e04/internal/age/ssh.go#L172
var curve25519P, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)

//EdPubToX25519 converts a ed25519 Public Key to a X25519 Public Key
func EdPubToX25519(pk ed25519.PublicKey) []byte {
	// ed25519.PublicKey is a little endian representation of the y-coordinate,
	// with the most significant bit set based on the sign of the x-coordinate.
	bigEndianY := make([]byte, ed25519.PublicKeySize)
	for i, b := range pk {
		bigEndianY[ed25519.PublicKeySize-i-1] = b
	}
	bigEndianY[0] &= 0b0111_1111

	// The Montgomery u-coordinate is derived through the bilinear map
	//
	//     u = (1 + y) / (1 - y)
	//
	// See https://blog.filippo.io/using-ed25519-keys-for-encryption.
	y := new(big.Int).SetBytes(bigEndianY)
	denom := big.NewInt(1)
	denom.ModInverse(denom.Sub(denom, y), curve25519P) // 1 / (1 - y)
	u := y.Mul(y.Add(y, big.NewInt(1)), denom)
	u.Mod(u, curve25519P)

	out := make([]byte, curve25519.PointSize)
	uBytes := u.Bytes()
	for i, b := range uBytes {
		out[len(uBytes)-i-1] = b
	}

	return out
}

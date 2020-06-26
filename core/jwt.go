package core

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/subtle"
	"errors"
	"github.com/vtolstov/jwt-go"
	"time"
)

type JwtClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}

// Generate a JWT token with the address and singed by the given private key
func GenerateJWTToken(addr HashAddress, key crypto.PrivateKey) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   addr.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(key)
}

// Validate a JWT token with the given public key and address
func ValidateJWTToken(tokenString string, addr HashAddress, key crypto.PublicKey) (*jwt.Token, error) {

	kf := func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, kf)
	if err != nil {
		return nil, err
	}

	// Make sure the token actually uses the correct signing method
	switch key.(type) {
	case *rsa.PrivateKey:
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("incorrect signing method")
		}
	case *ecdsa.PrivateKey:
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("incorrect signing method")
		}
	case ed25519.PrivateKey:
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("incorrect signing method")
		}
	default:
		return nil, errors.New("incorrect signing method")
	}

	// It should be a valid token
	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	// The standard claims should be valid
	err = token.Claims.Valid()
	if err != nil {
		return nil, err
	}

	// Check subject explicitly
	res := subtle.ConstantTimeCompare([]byte(token.Claims.(*jwt.StandardClaims).Subject), []byte(addr.String()))
	if res == 0 {
		return nil, errors.New("subject not valid")
	}

	return token, nil
}

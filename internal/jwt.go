package internal

import (
	"crypto/subtle"
	"errors"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	"github.com/vtolstov/jwt-go"
)

/*
 * @TODO
 * I don't really like this. Suppose we get access to a single JWT token. We can use the same token for every
 * single call in the next hour. Maybe we should limit each token for single-use (with a nonce that expires after
 * one hour, which means that expiresAt expires too), or maybe even limit the token for a single operation (add
 * request info?)
 */

const (
	invalidSigningMethod string = "invalid signing method"
)

// GenerateJWTToken generates a JWT token with the address and singed by the given private key
func GenerateJWTToken(addr hash.Hash, key bmcrypto.PrivKey) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: jwt.TimeFunc().Add(time.Hour * time.Duration(1)).Unix(),
		IssuedAt:  jwt.TimeFunc().Unix(),
		NotBefore: jwt.TimeFunc().Unix(),
		Subject:   addr.String(),
	}

	var signMethod jwt.SigningMethod
	switch key.Type {
	case bmcrypto.KeyTypeRSA:
		signMethod = jwt.SigningMethodRS256
	case bmcrypto.KeyTypeECDSA:
		signMethod = jwt.SigningMethodES256
	case bmcrypto.KeyTypeED25519:
		sm := &SigningMethodEdDSA{}
		signMethod = sm
		var edDSASigningMethod SigningMethodEdDSA
		jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })
	}

	token := jwt.NewWithClaims(signMethod, claims)

	return token.SignedString(key.K)
}

// ValidateJWTToken validates a JWT token with the given public key and address
func ValidateJWTToken(tokenString string, addr hash.Hash, key bmcrypto.PubKey) (*jwt.Token, error) {
	logrus.Tracef("validating JWT token: %s %s %s", tokenString, addr.String(), key.S)

	var edDSASigningMethod SigningMethodEdDSA
	jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })

	kf := func(token *jwt.Token) (interface{}, error) {
		/*
			// Make sure we signed with RS256
			if token.Method != jwt.SigningMethodRS256 {
				logrus.Trace("auth: jwt: " + invalidSigningMethod)
				return nil, errors.New(invalidSigningMethod)
			}
		*/
		return key.K, nil
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, kf)
	if err != nil {
		logrus.Trace("auth: jwt: ", err)
		return nil, err
	}

	// Make sure the token actually uses the correct signing method
	switch key.Type {
	case bmcrypto.KeyTypeRSA:
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			logrus.Tracef("auth: jwt: "+invalidSigningMethod)
			return nil, errors.New(invalidSigningMethod)
		}
	case bmcrypto.KeyTypeECDSA:
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			logrus.Tracef("auth: jwt: "+invalidSigningMethod)
			return nil, errors.New(invalidSigningMethod)
		}
	case bmcrypto.KeyTypeED25519:
		_, ok := token.Method.(*SigningMethodEdDSA)
		if !ok {
			logrus.Tracef("auth: jwt: "+invalidSigningMethod)
			return nil, errors.New(invalidSigningMethod)
		}
	default:
		logrus.Tracef("auth: jwt: "+invalidSigningMethod)
		return nil, errors.New(invalidSigningMethod)
	}

	// It should be a valid token
	if !token.Valid {
		logrus.Trace("auth: jwt: token not valid")
		logrus.Tracef("auth: jwt: %#v", token)
		return nil, errors.New("token not valid")
	}

	// The standard claims should be valid
	err = token.Claims.Valid()
	if err != nil {
		logrus.Trace("auth: jwt: ", err)
		return nil, err
	}

	// Check subject explicitly
	res := subtle.ConstantTimeCompare([]byte(token.Claims.(*jwt.StandardClaims).Subject), []byte(addr.String()))
	if res == 0 {
		logrus.Tracef("auth: jwt: subject does not match")
		return nil, errors.New("subject not valid")
	}

	return token, nil
}

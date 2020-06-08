package core

import (
    "crypto/subtle"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "github.com/dgrijalva/jwt-go"
    "time"
)

type JwtClaims struct {
    Address string `json:"address"`
    jwt.StandardClaims
}

// Generate a JWT token with the address and singed by the given private key
func GenerateJWTToken(addr HashAddress, privKey string) (string, error) {
    pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKey))
    if err != nil {
       return "", err
    }

    claims := &jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
        IssuedAt:  time.Now().Unix(),
        NotBefore: time.Now().Unix(),
        Subject:   addr.String(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

    return token.SignedString(pk)
}

// Validate a JWT token with the given public key and address
func ValidateJWTToken(tokenString string, addr HashAddress, pubKey string) (bool, error) {
    block, _ := pem.Decode([]byte(pubKey))
    pk, err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
        return false, err
    }

    kf := func(token *jwt.Token) (interface{}, error) {
        return pk, nil
    }

    claims := &jwt.StandardClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, kf)
    if err != nil {
        return false, err
    }

    // Make sure the token actually uses RSA for signing
    _, ok := token.Method.(*jwt.SigningMethodRSA)
    if ! ok {
        return false, errors.New("incorrect signing method")
    }

    // It should be a valid token
    if ! token.Valid {
        return false, errors.New("token not valid")
    }

    // The standard claims should be valid
    err = token.Claims.Valid()
    if err != nil {
        return false, err
    }

    // Check subject explicitly
    res := subtle.ConstantTimeCompare([]byte(token.Claims.(*jwt.StandardClaims).Subject), []byte(addr.String()))
    if res != 0 {
        return false, errors.New("subject not valid")
    }

    return token.Valid, nil
}

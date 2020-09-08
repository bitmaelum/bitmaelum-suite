package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)


func main() {
	logrus.SetLevel(logrus.TraceLevel)

	addr, err := address.NewHashFromHash(os.Args[1])
	if err != nil {
		panic(err)
	}

	key, err := bmcrypto.NewPubKey(os.Args[2])
	if err != nil {
		panic(err)
	}

	auth := "Bearer " + os.Args[3]

	if auth == "" {
		logrus.Trace("auth: empty auth string")
		os.Exit(1)
	}

	if len(auth) <= 6 || strings.ToUpper(auth[0:7]) != "BEARER " {
		logrus.Trace("auth: bearer not found")
		os.Exit(1)
	}
	tokenString := auth[7:]

	token, err := internal.ValidateJWTToken(tokenString, *addr, *key)
	if err == nil {
		fmt.Printf("Token validated correctly")
		spew.Dump(token)
		os.Exit(0)
	}

	logrus.Trace("auth: no key found that validates the token")
	os.Exit(1)
}

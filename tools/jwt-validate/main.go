// Copyright (c) 2022 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	addr, err := address.NewAddress(os.Args[1])
	if err != nil {
		panic(err)
	}

	key, err := bmcrypto.PublicKeyFromString(os.Args[2])
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

	token, err := api.ValidateJWTToken(tokenString, addr.Hash(), *key)
	if err == nil {
		fmt.Printf("Token validated correctly")
		spew.Dump(token)
		os.Exit(0)
	}

	logrus.Trace("auth: no key found that validates the token")
	os.Exit(1)
}

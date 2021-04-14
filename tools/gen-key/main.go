// Copyright (c) 2021 BitMaelum Authors
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

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <type>", os.Args[0])
		os.Exit(1)
	}

	keyType := os.Args[1]

	fmt.Printf("Generating key for %s...\n", keyType)
	kt, err := bmcrypto.FindKeyType(keyType)
	if err != nil {
		panic(err)
	}

	priv, pub, err := bmcrypto.GenerateKeyPair(kt);
	if err != nil {
		panic(err)
	}

	fmt.Println("Priv: ", priv.String());
	fmt.Println("Pub : ", pub.String());
}

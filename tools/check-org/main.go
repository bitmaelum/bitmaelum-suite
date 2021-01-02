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

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
)

type options struct {
	Config string `short:"c" long:"config" description:"Path to your configuration file"`
	Org    string `short:"o" long:"organisation" description:"Organisation" default:""`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	logrus.SetLevel(logrus.TraceLevel)

	v1, _ := organisation.NewValidationTypeFromString("dns bitmaelum.org")
	v2, _ := organisation.NewValidationTypeFromString("dns bitmaelum.com")
	v3, _ := organisation.NewValidationTypeFromString("dns evil-domain.xyz")

	o := organisation.Organisation{
		Hash:       hash.New("bitmaelum"),
		FullName:   "BitMaelum Org.",
		PublicKey:  bmcrypto.PubKey{},
		Validation: []organisation.ValidationType{*v1, *v2, *v3},
	}

	fmt.Printf("Organisation\n")
	fmt.Printf("  Hash: %s\n", o.Hash.String())
	fmt.Printf("  Validations: \n")
	for _, v := range o.Validation {
		if ok, err := v.Validate(o); err == nil && ok {
			fmt.Printf("    \U00002713 %s\n", v.String())
		} else {
			fmt.Printf("    \U00002717 %s\n", v.String())
		}
	}

}

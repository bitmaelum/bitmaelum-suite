package main

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
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
	a := address.NewHash("bitmaelum")

	o := organisation.Organisation{
		Addr:       a,
		Name:       "BitMaelum Org.",
		PublicKey:  bmcrypto.PubKey{},
		Validation: []organisation.ValidationType{*v1, *v2, *v3},
	}

	fmt.Printf("Organisation %s\n", o.Name)
	fmt.Printf("  Hash: %s\n", o.Addr.String())
	fmt.Printf("  Validations: \n")
	for _, v := range o.Validation {
		if ok, err := v.Validate(o); err == nil && ok {
			fmt.Printf("    \U00002713 %s\n", v.String())
		} else {
			fmt.Printf("    \U00002717 %s\n", v.String())
		}
	}

}

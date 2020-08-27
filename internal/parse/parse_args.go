package parse

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/xhit/go-str2duration/v2"
	"strconv"
	"strings"
	"time"
)

// ValidDuration gets a time duration string and return the time duration. Accepts single int as days
func ValidDuration(ds string) (time.Duration, error) {
	if ds == "" {
		return 0, nil
	}

	vd, err := str2duration.ParseDuration(ds)
	if err != nil {
		days, err := strconv.Atoi(ds)
		if err != nil {
			return 0, err
		}

		vd = time.Duration(days*24) * time.Hour
	}

	return vd, nil
}

// Permissions checks all permission and returns an error when a permission is not valid
func Permissions(perms []string) error {
	for _, p := range perms {
		p = strings.ToLower(p)

		found := false
		for _, ap := range apikey.AllPermissons {
			if p == ap {
				found = true
			}
		}
		if !found {
			return errors.New("unknown permission: " + p)
		}
	}

	return nil
}

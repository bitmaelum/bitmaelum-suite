package organisation

import (
	"net"
	"strings"
)

// Override for mocking purposes
var resolver DNSResolver = &DefaultResolver{}

func validateDNS(o Organisation, domain string) (bool, error) {
	oa := strings.ToLower(o.Hash.String())

	recs, err := resolver.LookupTXT("_bitmaelum." + domain)
	if err != nil {
		return false, err
	}

	for _, txt := range recs {
		if strings.ToLower(txt) == oa {
			return true, nil
		}
	}

	return false, nil
}

// DNSResolver is the interface for resolving DNS stuff
type DNSResolver interface {
	// We only need LookupTXT for now. Add more if they are needed.
	LookupTXT(host string) ([]string, error)
	SetCallbackTXT(callbackFunc)
}

// DefaultResolver is a resolver that will pass through to the net.Resolver
type DefaultResolver struct {
}

// SetCallbackTXT can set a callback for the LookupTXT resolver. Not used in the DefaultResolver
func (r *DefaultResolver) SetCallbackTXT(_ callbackFunc) {}

// LookupTXT passes through to the default net resolver
func (r *DefaultResolver) LookupTXT(host string) ([]string, error) {
	return net.LookupTXT(host)
}

type callbackFunc func() ([]string, error)

type mockResolver struct {
	callbackTxt callbackFunc
}

func (r *mockResolver) SetCallbackTXT(callback callbackFunc) {
	r.callbackTxt = callback
}

func (r *mockResolver) LookupTXT(name string) ([]string, error) {
	if r.callbackTxt != nil {
		return r.callbackTxt()
	}

	return []string{}, nil
}

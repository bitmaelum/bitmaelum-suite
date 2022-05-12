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
func (r *DefaultResolver) SetCallbackTXT(_ callbackFunc) {
	// Empty function
}

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

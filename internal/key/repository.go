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

package key

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

var (
	errKeyNotFound       = errors.New("key not found")
	errNeedsPointerValue = errors.New("needs a pointer value")
)

// GenericKey is a generic structure for storing and fetching keys
type GenericKey interface {
	GetID() string
	GetAddressHash() *hash.Hash
}

// StorageBackend is the main backend interface which can stores undefined structures. This can be api keys or auth
// key or even other kind of keys.
type StorageBackend interface {
	FetchByHash(h string, v interface{}) (interface{}, error)
	Fetch(ID string, v interface{}) error
	Store(v GenericKey) error
	Remove(v GenericKey) error
}

// AuthKeyRepo is the repository to the outside world. It allows us to fetch and store auth keys
type AuthKeyRepo interface {
	FetchByHash(h string) ([]AuthKeyType, error)
	Fetch(ID string) (*AuthKeyType, error)
	Store(v AuthKeyType) error
	Remove(v AuthKeyType) error
}

// APIKeyRepo is the repository to the outside world. It allows us to fetch and store api keys
type APIKeyRepo interface {
	FetchByHash(h string) ([]APIKeyType, error)
	Fetch(ID string) (*APIKeyType, error)
	Store(v APIKeyType) error
	Remove(v APIKeyType) error
}

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

package auth

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Lots of code is abstracted into functions. THis is to please sonarcloud duplication system
func TestAuthAPIKeyAuthenticate(t *testing.T) {
	var (
		exampleTestAddr = "example!"
		user1TestAddr   = "user-1!"
		user2TestAddr   = "user-2!"
		user3TestAddr   = "user-3!"
		expiredTestAddr = "expired!"
	)

	var apiKeyFixtures = []key.APIKeyType{
		key.NewAPIAccountKey(hash.New(user1TestAddr), []string{"a"}, time.Time{}, "my desc 1"),
		key.NewAPIAccountKey(hash.New(user2TestAddr), []string{"b"}, time.Time{}, "my desc 2"),
		key.NewAPIAccountKey(hash.New(user3TestAddr), []string{"b", "c"}, time.Time{}, "my desc 3"),
		key.NewAPIAccountKey(hash.New(expiredTestAddr), []string{"a", "b", "c"}, time.Unix(12521510, 0), "expired key"),
	}

	_, pubkey, err := testing2.ReadTestKey("../../../../testdata/key-ed25519-1.json")
	assert.NoError(t, err)

	container.Instance.SetShared("account", func() (interface{}, error) {
		return account.NewMockRepository(), nil
	})

	accountRepo := container.Instance.GetAccountRepo()
	// container.Set("account", func() interface{} { return accountRepo })
	_ = accountRepo.Create(hash.New(exampleTestAddr), *pubkey)
	_ = accountRepo.Create(hash.New(user1TestAddr), *pubkey)
	_ = accountRepo.Create(hash.New(user2TestAddr), *pubkey)
	_ = accountRepo.Create(hash.New(user3TestAddr), *pubkey)
	_ = accountRepo.Create(hash.New(expiredTestAddr), *pubkey)

	// 42 creates user1token
	rand.Seed(42)
	apiKeyRepo := key.NewAPIMockRepository()
	container.Instance.SetShared("api-key", func() (interface{}, error) { return apiKeyRepo, nil })
	for _, k := range apiKeyFixtures {
		// Create a new key, so it will randomize through our seed
		nk := key.NewAPIAccountKey(*k.AddressHash, k.Permissions, k.Expires, k.Desc)
		_ = apiKeyRepo.Store(nk)
	}

	a := APIKeyAuth{
		PermissionList: map[string][]string{
			"foo": {"a", "b", "c"},
			"bar": {"c", "a"},
			"baz": {"b"},
		},
	}

	var req *http.Request

	// No address
	req, _ = http.NewRequest("GET", "/foo", nil)
	checkFalse(t, &a, req)

	// No auth
	req = createReq("", exampleTestAddr)
	checkFalse(t, &a, req)

	// Address does not exist
	req = createReq("foobar", "doesnotexist!")
	checkFalse(t, &a, req)

	// ?
	req = createReq("foobar", exampleTestAddr)
	checkFalse(t, &a, req)

	// No key after bearer
	checkKey(t, a, false, "", "example!", "foo")

	// Expired key
	checkKey(t, a, false, "BMK-S7gYekwHUMGhWzGpld7aFPfYJK6SV75a", expiredTestAddr, "foo")

	user1token := "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io"

	// nonexisting route
	checkKey(t, a, false, user1token, user1TestAddr, "not-exist-in-perm-list")

	// Check all routes
	checkKey(t, a, false, user1token, user1TestAddr, "")    // no match
	checkKey(t, a, true, user1token, user1TestAddr, "foo")  // perm A
	checkKey(t, a, true, user1token, user1TestAddr, "bar")  // perm A
	checkKey(t, a, false, user1token, user1TestAddr, "baz") // no match

	// Token does not match for any other user
	checkKey(t, a, false, user1token, user2TestAddr, "foo")
	checkKey(t, a, false, user1token, user3TestAddr, "foo")

	user2token := "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi"

	checkKey(t, a, false, user2token, user1TestAddr, "")
	checkKey(t, a, false, user2token, user1TestAddr, "foo")
	checkKey(t, a, false, user2token, user1TestAddr, "bar")
	checkKey(t, a, false, user2token, user1TestAddr, "baz")

	checkKey(t, a, false, user2token, user2TestAddr, "")    // no match
	checkKey(t, a, true, user2token, user2TestAddr, "foo")  // perm B
	checkKey(t, a, false, user2token, user2TestAddr, "bar") // no match
	checkKey(t, a, true, user2token, user2TestAddr, "baz")  // perm B

	user3token := "BMK-FD4MY7O3gDk8Bg7W9LLxq2zGNO6q1Xh3"

	checkKey(t, a, false, user3token, user3TestAddr, "")   // no match
	checkKey(t, a, true, user3token, user3TestAddr, "foo") // Matches b and c
	checkKey(t, a, true, user3token, user3TestAddr, "bar") // Matches b
	checkKey(t, a, true, user3token, user3TestAddr, "baz") // Matches c
}

func createReq(auth string, addr string) *http.Request {
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Header.Set("authorization", auth)

	if addr != "" {
		req = mux.SetURLVars(req, map[string]string{
			"addr": hash.New(addr).String(),
		})
	}

	return req
}

func checkFalse(t *testing.T, a middleware.Authenticator, req *http.Request) {
	status, ctx, err := a.Authenticate(req, "")
	assert.Equal(t, status, middleware.AuthStatusPass)
	assert.Nil(t, ctx)
	assert.NoError(t, err)
}

func checkTrue(t *testing.T, a middleware.Authenticator, req *http.Request, hash string) {
	status, ctx, err := a.Authenticate(req, "")
	assert.Equal(t, status, middleware.AuthStatusSuccess)
	assert.Equal(t, hash, ctx.Value(middleware.AddressContext))
	assert.NoError(t, err)
}

func checkKey(t *testing.T, a APIKeyAuth, shouldPass bool, token, addr, routeName string) {
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Header.Set("authorization", "bearer "+token)
	req = mux.SetURLVars(req, map[string]string{
		"addr": hash.New(addr).String(),
	})

	status, ctx, err := a.Authenticate(req, routeName)
	if shouldPass {
		assert.Equal(t, status, middleware.AuthStatusSuccess)
		assert.NotNil(t, ctx)
		assert.NoError(t, err)
		// Check token in context
		k := ctx.Value(middleware.APIKeyContext).(*key.APIKeyType)
		assert.Equal(t, token, k.ID)
		assert.Equal(t, hash.New(addr).String(), k.AddressHash.String())
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, status, middleware.AuthStatusPass)
	assert.Nil(t, ctx)
}

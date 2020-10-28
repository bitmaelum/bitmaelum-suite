// Copyright (c) 2020 BitMaelum Authors
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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Lots of code is abstracted into functions. THis is to please sonarcloud duplication system

var apiKeyFixtures = []apikey.KeyType{
	apikey.NewAccountKey(hash.New("user-1!"), []string{"a"}, time.Time{}, "my desc 1"),
	apikey.NewAccountKey(hash.New("user-2!"), []string{"b"}, time.Time{}, "my desc 2"),
	apikey.NewAccountKey(hash.New("user-3!"), []string{"b", "c"}, time.Time{}, "my desc 3"),
	apikey.NewAccountKey(hash.New("expired!"), []string{"a", "b", "c"}, time.Unix(12521510, 0), "expired key"),
}

func TestAuthAPIKeyAuthenticate(t *testing.T) {
	t.SkipNow()

	_, pubkey, err := testing2.ReadTestKey("../../../../testdata/key-ed25519-1.json")
	assert.NoError(t, err)
	accountRepo = account.NewMockRepository
	_ = accountRepo().Create(hash.New("example!"), *pubkey)
	_ = accountRepo().Create(hash.New("user-1!"), *pubkey)
	_ = accountRepo().Create(hash.New("user-2!"), *pubkey)
	_ = accountRepo().Create(hash.New("user-3!"), *pubkey)
	_ = accountRepo().Create(hash.New("expired!"), *pubkey)

	// 42 creates BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io
	rand.Seed(42)
	apiKeyRepo = apikey.NewMockRepository
	for _, k := range apiKeyFixtures {
		// Create a new key, so it will randomize through our seed
		nk := apikey.NewAccountKey(*k.AddrHash, k.Permissions, k.Expires, k.Desc)
		_ = apiKeyRepo().Store(nk)
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
	req = createReq("", "example!")
	checkFalse(t, &a, req)

	// Address does not exist
	req = createReq("foobar", "doesnotexist!")
	checkFalse(t, &a, req)

	// ?
	req = createReq("foobar", "example!")
	checkFalse(t, &a, req)

	// No key after bearer
	checkKey(t, a, false, "", "example!", "foo")

	// Expired key
	checkKey(t, a, false, "BMK-S7gYekwHUMGhWzGpld7aFPfYJK6SV75a", "expired!", "foo")

	// nonexisting route
	checkKey(t, a, false, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-1!", "not-exist-in-perm-list")

	// Check all routes
	checkKey(t, a, false, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-1!", "")    // no match
	checkKey(t, a, true, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-1!", "foo")  // perm A
	checkKey(t, a, true, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-1!", "bar")  // perm A
	checkKey(t, a, false, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-1!", "baz") // no match

	// Token does not match for any other user
	checkKey(t, a, false, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-2!", "foo")
	checkKey(t, a, false, "BMK-dl2INvNSQTZ5zQu9MxNmGyAVmNkB33io", "user-3!", "foo")

	checkKey(t, a, false, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-1!", "")
	checkKey(t, a, false, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-1!", "foo")
	checkKey(t, a, false, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-1!", "bar")
	checkKey(t, a, false, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-1!", "baz")

	checkKey(t, a, false, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-2!", "")    // no match
	checkKey(t, a, true, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-2!", "foo")  // perm B
	checkKey(t, a, false, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-2!", "bar") // no match
	checkKey(t, a, true, "BMK-nwj2qrsh3xyC8OmCp1gObD0iOtQNQsLi", "user-2!", "baz")  // perm B

	checkKey(t, a, false, "BMK-FD4MY7O3gDk8Bg7W9LLxq2zGNO6q1Xh3", "user-3!", "")   // no match
	checkKey(t, a, true, "BMK-FD4MY7O3gDk8Bg7W9LLxq2zGNO6q1Xh3", "user-3!", "foo") // Matches b and c
	checkKey(t, a, true, "BMK-FD4MY7O3gDk8Bg7W9LLxq2zGNO6q1Xh3", "user-3!", "bar") // Matches b
	checkKey(t, a, true, "BMK-FD4MY7O3gDk8Bg7W9LLxq2zGNO6q1Xh3", "user-3!", "baz") // Matches c
}

func createReq(auth string, addr string) *http.Request {
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Header.Set("authorization", auth)

	if addr != "" {
		req = mux.SetURLVars(req, map[string]string{
			"addr": hash.New("example!").String(),
		})
	}

	return req
}

func checkFalse(t *testing.T, a middleware.Authenticator, req *http.Request) {
	ctx, ok := a.Authenticate(req, "")
	assert.False(t, ok)
	assert.Nil(t, ctx)
}

func checkKey(t *testing.T, a APIKeyAuth, pass bool, token, addr, routeName string) {
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Header.Set("authorization", "bearer "+token)
	req = mux.SetURLVars(req, map[string]string{
		"addr": hash.New(addr).String(),
	})

	ctx, ok := a.Authenticate(req, routeName)
	if pass {
		assert.True(t, ok)
		assert.NotNil(t, ctx)
		// Check token in context
		k := ctx.Value(APIKeyContext).(*apikey.KeyType)
		assert.Equal(t, token, k.ID)
		assert.Equal(t, hash.New(addr).String(), k.AddrHash.String())
		return
	}

	assert.False(t, ok)
	assert.Nil(t, ctx)
}

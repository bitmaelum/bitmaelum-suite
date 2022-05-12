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

package auth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/vtolstov/jwt-go"
)

// Lots of code is abstracted into functions. THis is to please sonarcloud duplication system

func TestAuthJwtAuthenticate(t *testing.T) {
	exampleTestAddr := "example!"

	_, pubkey, err := testing2.ReadTestKey("../../../../testdata/key-ed25519-1.json")
	assert.NoError(t, err)

	container.Instance.SetShared("account", func() (interface{}, error) {
		return account.NewMockRepository(), nil
	})

	accountRepo := container.Instance.GetAccountRepo()
	_ = accountRepo.Create(hash.New(exampleTestAddr), *pubkey)

	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	a := JwtAuth{}

	var (
		req    *http.Request
		ctx    context.Context
		status middleware.AuthStatus
	)

	// No address
	req, _ = http.NewRequest("GET", "/foo", nil)
	checkFalse(t, &a, req)

	// Address does not exist
	req = createReq("foobar", "doesnotexist!")
	checkFalse(t, &a, req)

	// No authorization
	req, _ = http.NewRequest("GET", "/foo", nil)
	req = mux.SetURLVars(req, map[string]string{
		"addr": hash.New(exampleTestAddr).String(),
	})
	checkFalse(t, &a, req)

	// No bearer key
	req = createReq("foobar", exampleTestAddr)
	checkFalse(t, &a, req)

	// Incorrect jwt token: not a token with the correct private key
	req, _ = http.NewRequest("GET", "/foo", nil)
	req.Header.Set("authorization", "bearer eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODU2OTYsImlhdCI6MTU3Nzg4MjA5NiwibmJmIjoxNTc3ODgyMDk2LCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.Bdm5brolKzTB4S-NQPTa93ubzPjejJb5hT8tpuRJG2Qpx3D0XrkAUAJNRyrQ2-aH188mfKmPcYeTXwd4qF3IAg")
	req = mux.SetURLVars(req, map[string]string{
		"addr": hash.New(exampleTestAddr).String(),
	})
	checkFalse(t, &a, req)

	// Correct JWT token
	req, _ = http.NewRequest("GET", "/foo", nil)
	req.Header.Set("authorization", "bearer eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODU2OTYsImlhdCI6MTU3Nzg4MjA5NiwibmJmIjoxNTc3ODgyMDk2LCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.EJmNoi18A0F_XGuel547ugFRcsIy3ZQj-NNp1JQB49zTdXHQ2Ob587CnYhUoREuHS-AJJAEHwuuAbsZIYkJoBw")
	req = mux.SetURLVars(req, map[string]string{
		"addr": hash.New(exampleTestAddr).String(),
	})

	status, ctx, err = a.Authenticate(req, "")
	assert.Equal(t, status, middleware.AuthStatusSuccess)
	assert.Equal(t, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab", ctx.Value(middleware.AddressContext))
	assert.NoError(t, err)

	// Incorrect time should return an explicit error
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 9, 9, 12, 34, 56, 0, time.UTC)
	}

	status, ctx, err = a.Authenticate(req, "")
	assert.Equal(t, status, middleware.AuthStatusFailure)
	assert.Nil(t, ctx)
	assert.EqualError(t, err, "token time not valid")

}

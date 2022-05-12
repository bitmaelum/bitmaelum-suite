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
	"net/http"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/vtolstov/jwt-go"
)

// Lots of code is abstracted into functions. THis is to please sonarcloud duplication system

const (
	// Generated on time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	token1 = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODIxMzAsImlhdCI6MTU3Nzg4MjA0MCwibmJmIjoxNTc3ODgyMDQwLCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.hZor6V8ZvjTto_jh4YiPceZJU9-0I6UmnZBNvrqV827i0p8f40qugSYu0mmA-JzMuHbLbR_jW-HccIgQBNAuAQ"
	token2 = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODIxMzAsImlhdCI6MTU3Nzg4MjA0MCwibmJmIjoxNTc3ODgyMDQwLCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.vK1YIFBN9u7DW8IJ9wWz95ZEJE7HZeGPLOkXwxWFiV5eZagQ5zpPZ8R3ah1P0EzW1q-d26Y8KQYrPXC05Xr1Bw"
	token3 = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODIxMzAsImlhdCI6MTU3Nzg4MjA0MCwibmJmIjoxNTc3ODgyMDQwLCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.x9wDsVSiB81PLpGvIbaBwj7img-tIdg7hwVE0Gldt-EOFgPK_aBA1njeU_0fbQlxVTVwoOCQaKJqGJ9xesJhAA"
	// Generated on return time.Date(2020, 01, 01, 13, 15, 41, 0, time.UTC)
	token4 = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODQ1OTAsImlhdCI6MTU3Nzg4NDUwMCwibmJmIjoxNTc3ODg0NTAwLCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.vRlguLgYtroXXv4ulOPs1qog3fdjRmLKG0Irg7j0qquWAff0mXU02l8cld7FqQCZVfUbeLlPSzVLaySgGS9KDw"
	token5 = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODQ1OTAsImlhdCI6MTU3Nzg4NDUwMCwibmJmIjoxNTc3ODg0NTAwLCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.OsgwNvYhOszxjH1PKt-zX-6TS5V5M6wjNNYX2ne4jStUp-FMD5nyNlSHKSJ6Gc9cllHLvgONBrqE2Y0kb3PmCg"
	token6 = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODQ1OTAsImlhdCI6MTU3Nzg4NDUwMCwibmJmIjoxNTc3ODg0NTAwLCJzdWIiOiIyZTQ1NTFkZTgwNGUyN2FhY2YyMGY5ZGY1YmUzZThjZDM4NGVkNjQ0ODhiMjFhYjA3OWZiNThlOGM5MDA2OGFiIn0.-u34k87bhnnoDaLhfeDuCAb_So0RkGaABt70HHHAMr5S9wABpt2OnQ2iSp5x9ihPz9dm5Cjmw33xlklOxX4TAQ"
)

func TestOnbehalfAuthJwtAuthenticate(t *testing.T) {
	privkey1, pubkey1, err := testing2.ReadTestKey("../../../../testdata/key-ed25519-1.json")
	assert.NoError(t, err)
	_, pubkey2, err := testing2.ReadTestKey("../../../../testdata/key-ed25519-2.json")
	assert.NoError(t, err)
	_, pubkey3, err := testing2.ReadTestKey("../../../../testdata/key-ed25519-3.json")
	assert.NoError(t, err)

	container.Instance.SetShared("account", func() (interface{}, error) {
		return account.NewMockRepository(), nil
	})
	container.Instance.SetShared("auth-key", func() (interface{}, error) {
		return key.NewAuthMockRepository(), nil
	})

	accountRepo := container.Instance.GetAccountRepo()
	_ = accountRepo.Create(hash.New("example!"), *pubkey1)

	// Sign key 2 with key 1
	ak := key.NewAuthKey(hash.New("example!"), pubkey2, "", time.Date(2020, 01, 01, 13, 00, 00, 0, time.UTC), "my desc")
	_ = ak.Sign(*privkey1)
	authRepo := container.Instance.GetAuthKeyRepo()
	_ = authRepo.Store(ak)

	// Sign key 3 with key 1
	ak = key.NewAuthKey(hash.New("example!"), pubkey3, "", time.Date(2020, 01, 01, 14, 00, 00, 0, time.UTC), "my desc")
	_ = ak.Sign(*privkey1)
	_ = authRepo.Store(ak)

	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	a := OnBehalfJwtAuth{}

	req, _ := http.NewRequest("GET", "/foo", nil)
	req = mux.SetURLVars(req, map[string]string{
		"addr": hash.New("example!").String(),
	})

	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	// Correct JWT token
	req.Header.Set("authorization", "bearer "+token1)
	checkFalse(t, &a, req)

	req.Header.Set("authorization", "bearer "+token2)
	checkTrue(t, &a, req, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab")

	req.Header.Set("authorization", "bearer "+token3)
	checkTrue(t, &a, req, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab")

	// Correct token, but expired token
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 13, 15, 41, 0, time.UTC)
	}

	req.Header.Set("authorization", "bearer "+token4)
	checkFalse(t, &a, req)

	// Just on the same time
	req.Header.Set("authorization", "bearer "+token5)
	checkTrue(t, &a, req, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab")

	req.Header.Set("authorization", "bearer "+token6)
	checkTrue(t, &a, req, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab")
}

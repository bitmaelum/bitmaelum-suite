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

package bmcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestKeyExchangeSet1 = []string{
	"rsa MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDTi87QrrSXLiILJ8/o8t48rS4XwtizoYbIfjV2K16ps9EmrNJvaQKW4MX5k6yBEocQTlxWq5YlYBM8V2NC8+W5TKzGT+MQuANZPiUNuZM/b3EgYFAciVAiudhqgEdT5vgt6BzKkymux3D0SJX1blH8X9tMyQnJoIPiT1xnML5P4gSG6SG9sSHNgrHrZwLjUr0yon9YGoDrCnYfvEn90ko+MK6k1tdf2XdBuPSvbRnjc3l5/uKaN+0SHaTCJA2e96qHq+mdiXUHcBkL5HxnAuKE/9BoEOWH6di7L8ndJF8zrIVhmflhB0cqZSPTwdYZXGdK4VkV4pMfDahDmYv6XKLpAgMBAAECggEAR57Ryj0bzwNDa1tzPH7dVtWbAVhqXYaWR1LTbsqIJhRG/z0LkcSPp905qaGhiaFoMNEW2hEFqGm6mXdMl+JTKEUZSZrKWWKzX4d2rArkG1nzhu6UsNScWOVqq8P6YiGUbJZlCQCB4DaNu2bHvmw3PaaGbJyzv5ukiv4rXpRWGlzhb3Iw//FGmbDAbQPnYrGLNXv0WDFlhME1gIZvOSytg6QDeCN3xYIJcmTJDCRuPemRN/MqCMSQ3qhqYLlqB+8vH902ozGeX8eO5ZiBJqc6oEjlKKwU45ER52QsXVmjsNCxf7L3IgN5CXNeuh/LKmFrOFWxjXHfOXqF7gwM91t2cQKBgQDp5YjADjGm88UyUicTYu9Ta2o2gxfMiyPb73IhfPJ/ijqlUSCUN9Ji/6Vyrd5/OjqLENWKUsjrMMFFhqOmLem27yuRNAGUQeoO4fibf2oNSNMxvsEFfwlVJAq4ixPJDNyOne2vWEwb9u7w4YK+ssr3cOZL+mG2xk/5XcK36g5S3wKBgQDniZEcqkGzTiQFjgf+7qHtCwaDNiNzvQCPL0ZsQfiur4NQrL/osGIZ4jHV5eU/k/LuHcrZx5x+UAQqVISScFkkew2VvgfKaxdie8zhJSh8ic5RBGZg/B3gw5xHWIyJSxUp8PztuH2Raj0ECh+lzDtwYpPKGyQLDxwdeB1/Pq7LNwKBgG+juHb7D2YBuqD/J1mQgm0Nux+TyNs/mnkSvCYRzmlj4AQiSeuVDV1lamHnbWjKsUDJYzNnujDQD6AQ2LGr/n7rf58J9KsAHyjFYPVPhp4aoXt/8f+emCTEVD2rXGE9O1TzOozUF1fNsFTXPqGpE0mx4KppMxSbaXa78wH3vKh/AoGBAJf4645tEgKm323l88mYyB/WhMfK2So2fA9/cDHOe3PtL7vcJ3qLi1iB50QGSZqZeXJhi6u2ITmnO5StNPcJVvli61/GA0cRU6AIskl1IkXcDdePk8NEuDe3LPSHYncbGSEWVG2UEpdHrBTisDMbAkiZ63dUqSu5FzMgi/vhIMmxAoGBAOHgQ7xMwPOw7mpczSdT8gUqNwGFWKIVMo0YJFW20WntJgr/aFY70XoW4ENIHJyZFEl2GSuHiI6S8d2yLJdbSmjlj06wgwCsEk2Bc8e81BSVHisH8e2kSApyjMOxjgfUZoWcEHrWJeIcZykowiPrF39HDMbOtwv94n96+v90T4Nq",
	"rsa MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA04vO0K60ly4iCyfP6PLePK0uF8LYs6GGyH41diteqbPRJqzSb2kCluDF+ZOsgRKHEE5cVquWJWATPFdjQvPluUysxk/jELgDWT4lDbmTP29xIGBQHIlQIrnYaoBHU+b4LegcypMprsdw9EiV9W5R/F/bTMkJyaCD4k9cZzC+T+IEhukhvbEhzYKx62cC41K9MqJ/WBqA6wp2H7xJ/dJKPjCupNbXX9l3Qbj0r20Z43N5ef7imjftEh2kwiQNnveqh6vpnYl1B3AZC+R8ZwLihP/QaBDlh+nYuy/J3SRfM6yFYZn5YQdHKmUj08HWGVxnSuFZFeKTHw2oQ5mL+lyi6QIDAQAB",
}

var TestKeyExchangeSet2 = []string{
	"ed25519 MC4CAQAwBQYDK2VwBCIEIJcq+oL3JAroQaJ63+iuAZ41s8/fwLmkq5CuSMK/tqyr",
	"ed25519 MCowBQYDK2VwAyEAFLbQmBHt0si4scDb6lLiLAnD9V2p8Jw2tqwsg1aoYRQ=",
}
var TestKeyExchangeSet3 = []string{
	"ed25519 MC4CAQAwBQYDK2VwBCIEIP7FYQY5CPyMgBCKN5M95CC+T+leCyRwkQBZNDtmQwtr",
	"ed25519 MCowBQYDK2VwAyEA8If1INF+kc3n3wa/PplsNRQCylsxEJq70fSqzCHV2EU=",
}
var TestKeyExchangeSet4 = []string{
	"ecdsa MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDAzck3f6PLACzrqcH4i2XkDcOUwTMd6uOGH+B9BIIh1A2vjvF/OI8PsclbQ1emFW8qhZANiAASTZ4OiS7ntWsA7Yg7A/ZxhkYa8rrPcbbwfsBkrnQRIJAnaKZp70kt27SghVQAXzPrl4BoqlvR/B5LlSAb5/R/+EiERyWoMyOe9gi+CCItLvyyhj/shAWiHm1WiOMbWQko=",
	"ecdsa MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEk2eDoku57VrAO2IOwP2cYZGGvK6z3G28H7AZK50ESCQJ2imae9JLdu0oIVUAF8z65eAaKpb0fweS5UgG+f0f/hIhEclqDMjnvYIvggiLS78soY/7IQFoh5tVojjG1kJK",
}
var TestKeyExchangeSet5 = []string{
	"ecdsa MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDBN/DY/gyLxwaIUoULm+ONL90PFKJmCIpZ8cKKzGVeQa8WYSlbF19yeO7PxTeVcYNChZANiAAQIglBvq9D4ScAT1GlJT13t1RUrjly42dAbjatyW65yu9EQm06y8RdA45DyvCQzLoet8XKpYSUQrjyBfK3mhFn9uHr6VTeqIluBlXWofm7D5fT7o/p8udyUKREERE+G9fs=",
	"ecdsa MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAECIJQb6vQ+EnAE9RpSU9d7dUVK45cuNnQG42rcluucrvREJtOsvEXQOOQ8rwkMy6HrfFyqWElEK48gXyt5oRZ/bh6+lU3qiJbgZV1qH5uw+X0+6P6fLnclCkRBERPhvX7",
}

var (
	ecdsaSharedSecret   = []byte{198, 119, 229, 119, 133, 42, 232, 116, 170, 231, 202, 204, 162, 20, 175, 170, 56, 164, 144, 5, 72, 186, 178, 228, 237, 180, 249, 250, 209, 92, 49, 234, 26, 82, 65, 174, 174, 193, 178, 177, 211, 105, 215, 179, 21, 173, 245, 199}
	ed25519SharedSecret = []byte{145, 48, 204, 130, 242, 13, 147, 255, 174, 10, 185, 101, 47, 22, 31, 169, 98, 60, 176, 85, 18, 98, 178, 167, 181, 219, 74, 63, 176, 57, 13, 122}
)

func TestKeyExchange(t *testing.T) {
	// Incorrect key type
	privAliceK, _ := PrivKeyFromString(TestKeyExchangeSet1[0])
	pubBobK, _ := PubKeyFromString(TestKeyExchangeSet1[1])
	_, err := KeyExchange(*privAliceK, *pubBobK)
	assert.Error(t, err)

	// ecdsa key exchange (ecdh) alice->bob
	privAliceK, _ = PrivKeyFromString(TestKeyExchangeSet4[0])
	pubBobK, _ = PubKeyFromString(TestKeyExchangeSet5[1])
	k, err := KeyExchange(*privAliceK, *pubBobK)
	assert.NoError(t, err)
	assert.Equal(t, ecdsaSharedSecret, k)

	// ecdsa key exchange (ecdh) bob->alice
	privBobK, _ := PrivKeyFromString(TestKeyExchangeSet5[0])
	pubAliceK, _ := PubKeyFromString(TestKeyExchangeSet4[1])
	k, err = KeyExchange(*privBobK, *pubAliceK)
	assert.NoError(t, err)
	assert.Equal(t, ecdsaSharedSecret, k)

	// ed25519 key exchange (ecdh on curve 25519) alice->bob
	privAliceK, _ = PrivKeyFromString(TestKeyExchangeSet2[0])
	pubBobK, _ = PubKeyFromString(TestKeyExchangeSet3[1])
	k, err = KeyExchange(*privAliceK, *pubBobK)
	assert.NoError(t, err)
	assert.Equal(t, ed25519SharedSecret, k)

	// ed25519 key exchange (ecdh on curve 25519) bob->alice
	privBobK, _ = PrivKeyFromString(TestKeyExchangeSet3[0])
	pubAliceK, _ = PubKeyFromString(TestKeyExchangeSet2[1])
	k, err = KeyExchange(*privBobK, *pubAliceK)
	assert.NoError(t, err)
	assert.Equal(t, ed25519SharedSecret, k)

}

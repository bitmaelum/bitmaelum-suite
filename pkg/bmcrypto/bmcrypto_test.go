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

package bmcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKeyType(t *testing.T) {
	kt, err := FindKeyType("rsa")
	assert.NoError(t, err)
	assert.Equal(t, "rsa", kt.String())

	kt, err = FindKeyType("does-not-exist")
	assert.Error(t, err)
	assert.Nil(t, kt)
}

func TestGenerateKeyPair(t *testing.T) {
	_, _, err := GenerateKeyPair(nil)
	assert.Error(t, err)
}

func TestKeyExchangeError(t *testing.T) {
	privKey, err := PrivateKeyFromString("rsa MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDTi87QrrSXLiILJ8/o8t48rS4XwtizoYbIfjV2K16ps9EmrNJvaQKW4MX5k6yBEocQTlxWq5YlYBM8V2NC8+W5TKzGT+MQuANZPiUNuZM/b3EgYFAciVAiudhqgEdT5vgt6BzKkymux3D0SJX1blH8X9tMyQnJoIPiT1xnML5P4gSG6SG9sSHNgrHrZwLjUr0yon9YGoDrCnYfvEn90ko+MK6k1tdf2XdBuPSvbRnjc3l5/uKaN+0SHaTCJA2e96qHq+mdiXUHcBkL5HxnAuKE/9BoEOWH6di7L8ndJF8zrIVhmflhB0cqZSPTwdYZXGdK4VkV4pMfDahDmYv6XKLpAgMBAAECggEAR57Ryj0bzwNDa1tzPH7dVtWbAVhqXYaWR1LTbsqIJhRG/z0LkcSPp905qaGhiaFoMNEW2hEFqGm6mXdMl+JTKEUZSZrKWWKzX4d2rArkG1nzhu6UsNScWOVqq8P6YiGUbJZlCQCB4DaNu2bHvmw3PaaGbJyzv5ukiv4rXpRWGlzhb3Iw//FGmbDAbQPnYrGLNXv0WDFlhME1gIZvOSytg6QDeCN3xYIJcmTJDCRuPemRN/MqCMSQ3qhqYLlqB+8vH902ozGeX8eO5ZiBJqc6oEjlKKwU45ER52QsXVmjsNCxf7L3IgN5CXNeuh/LKmFrOFWxjXHfOXqF7gwM91t2cQKBgQDp5YjADjGm88UyUicTYu9Ta2o2gxfMiyPb73IhfPJ/ijqlUSCUN9Ji/6Vyrd5/OjqLENWKUsjrMMFFhqOmLem27yuRNAGUQeoO4fibf2oNSNMxvsEFfwlVJAq4ixPJDNyOne2vWEwb9u7w4YK+ssr3cOZL+mG2xk/5XcK36g5S3wKBgQDniZEcqkGzTiQFjgf+7qHtCwaDNiNzvQCPL0ZsQfiur4NQrL/osGIZ4jHV5eU/k/LuHcrZx5x+UAQqVISScFkkew2VvgfKaxdie8zhJSh8ic5RBGZg/B3gw5xHWIyJSxUp8PztuH2Raj0ECh+lzDtwYpPKGyQLDxwdeB1/Pq7LNwKBgG+juHb7D2YBuqD/J1mQgm0Nux+TyNs/mnkSvCYRzmlj4AQiSeuVDV1lamHnbWjKsUDJYzNnujDQD6AQ2LGr/n7rf58J9KsAHyjFYPVPhp4aoXt/8f+emCTEVD2rXGE9O1TzOozUF1fNsFTXPqGpE0mx4KppMxSbaXa78wH3vKh/AoGBAJf4645tEgKm323l88mYyB/WhMfK2So2fA9/cDHOe3PtL7vcJ3qLi1iB50QGSZqZeXJhi6u2ITmnO5StNPcJVvli61/GA0cRU6AIskl1IkXcDdePk8NEuDe3LPSHYncbGSEWVG2UEpdHrBTisDMbAkiZ63dUqSu5FzMgi/vhIMmxAoGBAOHgQ7xMwPOw7mpczSdT8gUqNwGFWKIVMo0YJFW20WntJgr/aFY70XoW4ENIHJyZFEl2GSuHiI6S8d2yLJdbSmjlj06wgwCsEk2Bc8e81BSVHisH8e2kSApyjMOxjgfUZoWcEHrWJeIcZykowiPrF39HDMbOtwv94n96+v90T4Nq")
	assert.NoError(t, err)

	pubKey, err := PublicKeyFromString("rsa MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA04vO0K60ly4iCyfP6PLePK0uF8LYs6GGyH41diteqbPRJqzSb2kCluDF+ZOsgRKHEE5cVquWJWATPFdjQvPluUysxk/jELgDWT4lDbmTP29xIGBQHIlQIrnYaoBHU+b4LegcypMprsdw9EiV9W5R/F/bTMkJyaCD4k9cZzC+T+IEhukhvbEhzYKx62cC41K9MqJ/WBqA6wp2H7xJ/dJKPjCupNbXX9l3Qbj0r20Z43N5ef7imjftEh2kwiQNnveqh6vpnYl1B3AZC+R8ZwLihP/QaBDlh+nYuy/J3SRfM6yFYZn5YQdHKmUj08HWGVxnSuFZFeKTHw2oQ5mL+lyi6QIDAQAB")
	assert.NoError(t, err)

	_, err = KeyExchange(*privKey, *pubKey)
	assert.Error(t, err)
}

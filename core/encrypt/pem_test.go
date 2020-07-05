package encrypt

import (
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestPEM(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/mykey.pub")
	pubKey, err := PEMToPubKey(data)
	assert.NoError(t, err)
	assert.IsType(t, (*rsa.PublicKey)(nil), pubKey)

	pem, err := PubKeyToPEM(pubKey)
	assert.NoError(t, err)
	assert.Equal(t, "-----BEGIN RSA PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8\npvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5a\nG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxX\nPOhC1QUcllt08yBB3wIDAQAB\n-----END RSA PUBLIC KEY-----\n", pem)

	data, _ = ioutil.ReadFile("./testdata/mykey.pem")
	privKey, _ := PEMToPrivKey(data)
	assert.NoError(t, err)
	assert.IsType(t, (*rsa.PrivateKey)(nil), privKey)

	pem, err = PrivKeyToPEM(privKey)
	assert.NoError(t, err)
	assert.Equal(t, "-----BEGIN PRIVATE KEY-----\nMIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBALnuoL8F6hhwzqKN\nrO5oJ0mRtPym8+kUQNXN/1lD1rGzBdTebLBI5pHfB+WuJKOa9F9Q2z8ECOyGC0mq\nC+JRfHuUTlob2ZxTK+ibIUWl8JHaXYu+jDQfzLsmgvrLExJB/yCkBxpDM8Hx2x8Z\n74OoD7nCPFc86ELVBRyWW3TzIEHfAgMBAAECgYEApJYMqyu0HnB1KcWZx+xgoqod\neOzcyn0IK3rPR5haizB6wAUoVyAhIg04s2LkwgJfwaQUgAK1V5IMmeex32PceREP\nuTDGa9VVeW1oitusS6cqB1ErUAKTYGRvtRQMiYOOJnO0m8JIv5Yu70WmGfwh8/PA\nAg9YAooGgPWS4jrJ8UECQQDp0mmhMwoOUIQx4lZlVSeEvsr84GmFYs7txl9rGjAF\n4YZe4wFaJrJKx14gNXtL88OegcErAiRKxJ+PNZim26rhAkEAy5FfCipkco8eMxML\nLPcLufAULOCaYozk9F5lTSzj+10W79UY5KdUNLMLkkVxjjwKzTdMHBnbA+G//evh\nrthEvwJBALF51FNWujtDQhPbCFjB2c0YRFrMu0tTRF2WRLa2mdzc4XEEPPKAjLPV\nv8wSzBNKYyDcvBI4/fMCa1n4BHYiJgECQEQ2CeqWGeJpIm1qzCvc/Ajp8Vi4SgML\nwww/NSKPEDWCE9V5SkYT65tdA1uX+Vz3sYbiKJah80lxs6uaBp1XynECQGGyfZP0\nNKgN14iZF5zGzf9d6PgRhCmP/dfjfpGsWe2NPTrAPrAX2zkK8+bP9+1QEnvUhWUF\nBR3eM6BA9HMSE0I=\n-----END PRIVATE KEY-----\n", pem)
}

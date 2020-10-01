package organisation

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// Override for testing purposes
var timeNow = time.Now

// GenerateInviteToken generates an invite token for the given
func GenerateInviteToken(addr *address.Address, routingID string, validUntil time.Time, key bmcrypto.PrivKey) (string, error) {
	ts := strconv.FormatInt(validUntil.Unix(), 10)

	// Not an organisation hash, so no token is needed
	if !addr.IsOrganisationAddress() {
		return "", errors.New("not an organisation address")
	}

	dataToSign := addr.Hash().String() + routingID + ts
	hash := sha256.Sum256([]byte(dataToSign))
	signedData, err := bmcrypto.Sign(key, hash[:])
	if err != nil {
		return "", errors.New("could not sign data")
	}

	token := addr.Hash().String() + ":" + routingID + ":" + ts + ":" + string(signedData)

	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

func VerifyInviteToken(token string, addr *address.Address, routingID string, key bmcrypto.PubKey) bool {
	tokenData, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(tokenData), ":", 4)
	if len(parts) != 4 {
		return false
	}

	// Check signature first
	hash := sha256.Sum256([]byte(parts[0] + parts[1] + parts[2]))
	ok, err := bmcrypto.Verify(key, hash[:], []byte(parts[3]))
	if err != nil || !ok {
		return false
	}

	// Check address
	if addr.Hash().String() != parts[0] {
		return false
	}

	// Check routing
	if routingID != parts[1] {
		return false
	}

	// Check expiry
	ts, err := strconv.Atoi(parts[2])
	if err != nil {
		return false
	}
	expiry := time.Unix(int64(ts), 0)
	return !timeNow().After(expiry)
}

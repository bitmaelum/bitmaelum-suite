package subscription

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// Subscription is a tuple that can be used to identify a mailing list user. If one is found, we are allowed to skip
// proof-of-work during uploading messages.
type Subscription struct {
	From           address.Hash
	To             address.Hash
	SubscriptionID string
}

// New returns a new subscription
func New(from, to address.Hash, subscriptionID string) Subscription {
	return Subscription{
		From:           from,
		To:             to,
		SubscriptionID: subscriptionID,
	}
}

// Repository is the interface that needs to be implemented by subscription storage
type Repository interface {
	Has(sub *Subscription) bool
	Store(sub *Subscription) error
	Remove(sub *Subscription) error
}

// Generate a key that can be used for reading / writing the subscription info
func createKey(sub *Subscription) string {
	data := fmt.Sprintf("%s-%s-%s", sub.From.String(), sub.To.String(), sub.SubscriptionID)
	h := sha256.Sum256([]byte(data))

	return "sub-" + hex.EncodeToString(h[:])
}

package subscription

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// Subscription is a tuple that can be used to identify a mailing list user. If one is found, we are allowed to skip
// proof-of-work during uploading messages.
type Subscription struct {
	From           address.HashAddress
	To             address.HashAddress
	SubscriptionID string
}

// New returns a new subscription
func New(from, to address.HashAddress, subscriptionID string) Subscription {
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

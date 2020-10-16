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

package subscription

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// Subscription is a tuple that can be used to identify a mailing list user. If one is found, we are allowed to skip
// proof-of-work during uploading messages.
type Subscription struct {
	From           hash.Hash
	To             hash.Hash
	SubscriptionID string
}

// New returns a new subscription
func New(from, to hash.Hash, subscriptionID string) Subscription {
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

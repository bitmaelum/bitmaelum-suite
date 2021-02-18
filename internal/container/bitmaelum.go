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

package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/store"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
)

const (
	// APIKey key
	APIKey = "api-key"

	// AuthKey key
	AuthKey = "auth-key"

	// Resolver key
	Resolver = "resolver"

	// Subscription key
	Subscription = "subscription"

	// Ticket key
	Ticket = "ticket"

	// Webhook key
	Webhook = "webhook"
)

// Instance is the main bitmaelum service container
var Instance = Type{
	definitions: make(map[string]*ServiceDefinition),
	resolved:    make(map[string]interface{}),
}

// GetAPIKeyRepo will return the current api key repository
func (c *Type) GetAPIKeyRepo() key.APIKeyRepo {
	retValue, _ := c.Get(APIKey).(key.APIKeyRepo)
	return retValue
}

// GetAuthKeyRepo will return the current auth key repository
func (c *Type) GetAuthKeyRepo() key.AuthKeyRepo {
	retValue, _ := c.Get(AuthKey).(key.AuthKeyRepo)
	return retValue
}

// GetResolveService will return the current resolver service
func (c *Type) GetResolveService() *resolver.Service {
	retValue, _ := c.Get(Resolver).(*resolver.Service)
	return retValue
}

// GetSubscriptionRepo will return the current subscription repository
func (c *Type) GetSubscriptionRepo() subscription.Repository {
	retValue, _ := c.Get(Subscription).(subscription.Repository)
	return retValue
}

// GetTicketRepo will return the current ticket repository
func (c *Type) GetTicketRepo() ticket.Repository {
	retValue, _ := c.Get(Ticket).(ticket.Repository)
	return retValue
}

// GetWebhookRepo will return the current webhook repository
func (c *Type) GetWebhookRepo() webhook.Repository {
	retValue, _ := c.Get(Webhook).(webhook.Repository)
	return retValue
}

// GetStoreRepo will return the current store repository
func (c *Type) GetStoreRepo() store.Repository {
	return c.Get("store").(store.Repository)
}

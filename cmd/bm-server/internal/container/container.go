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

package container

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/account"
	maincontainer "github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
)

/*
 * MultiContainer is a container struct that acts like a switch between the "client" container, and the "general" container.
 * The reason this is not a single container, is that the client container has additional methods which can be called, for
 * instance: GetAccountRepo().
 *
 * So, all the code in the bm-server directory will use this container, so it can fetch the Account Repo when needed, but
 * also can fetch other repositories from the same instance. From the code point of view, all is located in the same container.
 *
 * I really think there are better ways to deal with this, but this is the best i can do for now.
 */

// MultiContainer is a struct that holds both the general and client container
type MultiContainer struct {
	general maincontainer.Container // The general container used in the internal packages
	client  maincontainer.Container // The client container used in this command
}

// Instance is the main bitmaelum service container
var Instance = MultiContainer{
	general: &maincontainer.Instance,
	client:  maincontainer.NewContainer(),
}

// SetShared will set a definition inside the client container
func (c *MultiContainer) SetShared(key string, f maincontainer.ServiceFunc) {
	c.client.SetShared(key, f)
}

// SetNonShared will set a definition inside the client container
func (c *MultiContainer) SetNonShared(key string, f maincontainer.ServiceFunc) {
	c.client.SetNonShared(key, f)
}

// Get will fetch a definition from the client container
func (c *MultiContainer) Get(key string) interface{} {
	return c.client.Get(key)
}

// GetAccountRepo will return the current account repository
func (c *MultiContainer) GetAccountRepo() account.Repository {
	return c.client.Get("account").(account.Repository)
}

// GetAPIKeyRepo will return the current api key repository
func (c *MultiContainer) GetAPIKeyRepo() key.APIKeyRepo {
	if c.client.Has("api-key") {
		return c.client.Get("api-key").(key.APIKeyRepo)
	}

	return c.general.Get("api-key").(key.APIKeyRepo)
}

// GetAuthKeyRepo will return the current auth key repository
func (c *MultiContainer) GetAuthKeyRepo() key.AuthKeyRepo {
	if c.client.Has("auth-key") {
		return c.client.Get("auth-key").(key.AuthKeyRepo)
	}

	return c.general.Get("auth-key").(key.AuthKeyRepo)
}

// GetResolveService will return the current resolver service
func (c *MultiContainer) GetResolveService() *resolver.Service {
	if c.client.Has("resolver") {
		return c.client.Get("resolver").(*resolver.Service)
	}

	return c.general.Get("resolver").(*resolver.Service)
}

// GetSubscriptionRepo will return the current subscription repository
func (c *MultiContainer) GetSubscriptionRepo() subscription.Repository {
	if c.client.Has("subscription") {
		return c.client.Get("subscription").(subscription.Repository)
	}

	return c.general.Get("subscription").(subscription.Repository)
}

// GetTicketRepo will return the current ticket repository
func (c *MultiContainer) GetTicketRepo() ticket.Repository {
	if c.client.Has("ticket") {
		return c.client.Get("ticket").(ticket.Repository)
	}

	return c.general.Get("ticket").(ticket.Repository)
}

// GetWebhookRepo will return the current web hook repository
func (c *MultiContainer) GetWebhookRepo() webhook.Repository {
	if c.client.Has("webhook") {
		return c.client.Get("webhook").(webhook.Repository)
	}

	return c.general.Get("webhook").(webhook.Repository)
}

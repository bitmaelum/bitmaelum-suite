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

package container

import (
	maincontainer "github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
)

/*
 * For more information about MultiContainers, please see:
 *    github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container/container.go
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

// GetAPIKeyRepo will return the current api key repository
func (c *MultiContainer) GetAPIKeyRepo() key.APIKeyRepo {
	if c.client.Has(maincontainer.APIKey) {
		return c.client.Get(maincontainer.APIKey).(key.APIKeyRepo)
	}

	return c.general.Get(maincontainer.APIKey).(key.APIKeyRepo)
}

// GetAuthKeyRepo will return the current auth key repository
func (c *MultiContainer) GetAuthKeyRepo() key.AuthKeyRepo {
	if c.client.Has(maincontainer.AuthKey) {
		return c.client.Get(maincontainer.AuthKey).(key.AuthKeyRepo)
	}

	return c.general.Get(maincontainer.AuthKey).(key.AuthKeyRepo)
}

// GetResolveService will return the current resolver service
func (c *MultiContainer) GetResolveService() *resolver.Service {
	if c.client.Has(maincontainer.Resolver) {
		return c.client.Get(maincontainer.Resolver).(*resolver.Service)
	}

	return c.general.Get(maincontainer.Resolver).(*resolver.Service)
}

// GetSubscriptionRepo will return the current subscription repository
func (c *MultiContainer) GetSubscriptionRepo() subscription.Repository {
	if c.client.Has(maincontainer.Subscription) {
		return c.client.Get(maincontainer.Subscription).(subscription.Repository)
	}

	return c.general.Get(maincontainer.Subscription).(subscription.Repository)
}

// GetTicketRepo will return the current ticket repository
func (c *MultiContainer) GetTicketRepo() ticket.Repository {
	if c.client.Has(maincontainer.Ticket) {
		return c.client.Get(maincontainer.Ticket).(ticket.Repository)
	}

	return c.general.Get(maincontainer.Ticket).(ticket.Repository)
}

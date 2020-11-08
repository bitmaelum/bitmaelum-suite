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
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
)

// Instance is the main bitmaelum service container
var Instance = Container{
	definitions: make(map[string]*ServiceDefinition),
	resolved:    make(map[string]interface{}),
}

// GetAPIKeyRepo will return the current api key repository
func (c *Container) GetAPIKeyRepo() key.APIKeyRepo {
	return c.Get("api-key").(key.APIKeyRepo)
}

// GetAuthKeyRepo will return the current auth key repository
func (c *Container) GetAuthKeyRepo() key.AuthKeyRepo {
	return c.Get("auth-key").(key.AuthKeyRepo)
}

// GetResolveService will return the current resolver service
func (c *Container) GetResolveService() *resolver.Service {
	return c.Get("resolver").(*resolver.Service)
}

// GetSubscriptionRepo will return the current subscription repository
func (c *Container) GetSubscriptionRepo() subscription.Repository {
	return c.Get("subscription").(subscription.Repository)
}

// GetTicketRepo will return the current ticket repository
func (c *Container) GetTicketRepo() ticket.Repository {
	return c.Get("ticket").(ticket.Repository)
}


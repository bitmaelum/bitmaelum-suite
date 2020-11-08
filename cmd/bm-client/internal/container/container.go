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
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
)

type ClientContainer struct {
	Container container.Container
	Instance *container.Container
}

var Instance = ClientContainer{
	Container: container.NewContainer(),
	Instance: &container.Instance,
}

func (c *ClientContainer) SetShared(key string, f container.ServiceFunc) {
	c.Container.SetShared(key, f)
}

// GetAPIKeyRepo will return the current api key repository
func (c *ClientContainer) GetAPIKeyRepo() key.APIKeyRepo {
	return c.Instance.Get("api-key").(key.APIKeyRepo)
}

// GetAuthKeyRepo will return the current auth key repository
func (c *ClientContainer) GetAuthKeyRepo() key.AuthKeyRepo {
	return c.Instance.Get("auth-key").(key.AuthKeyRepo)
}

// GetResolveService will return the current resolver service
func (c *ClientContainer) GetResolveService() *resolver.Service {
	return c.Instance.Get("resolver").(*resolver.Service)
}

// GetSubscriptionRepo will return the current subscription repository
func (c *ClientContainer) GetSubscriptionRepo() subscription.Repository {
	return c.Instance.Get("subscription").(subscription.Repository)
}

// GetTicketRepo will return the current ticket repository
func (c *ClientContainer) GetTicketRepo() ticket.Repository {
	return c.Instance.Get("ticket").(ticket.Repository)
}


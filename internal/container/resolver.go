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
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
)

// We can have multiple resolvers to resolve a single address. We could resolve locally, remotely through
// resolver-services, or through DHT. We chain them all together with the ChainRepository

var (
	resolveOnce    sync.Once
	resolveService *resolver.Service
)

func setupResolverService() (interface{}, error) {
	resolveOnce.Do(func() {
		repo := resolver.NewChainRepository()
		if config.Server.Resolver.Sqlite.Enabled {
			r, err := getSQLiteRepository(config.Server.Resolver.Sqlite.Dsn)
			if err == nil {
				_ = repo.Add(r)
			}
		}

		// We add either the client or the server resolver
		if config.Client.Resolver.Remote.Enabled {
			_ = repo.Add(*getRemoteRepository(config.Client.Resolver.Remote.URL, config.Client.Server.DebugHTTP, config.Client.Resolver.Remote.AllowInsecure))
		}
		if config.Server.Resolver.Remote.Enabled {
			_ = repo.Add(*getRemoteRepository(config.Server.Resolver.Remote.URL, false, config.Client.Resolver.Remote.AllowInsecure))
		}

		resolveService = resolver.KeyRetrievalService(repo)
	})

	return resolveService, nil
}

func getRemoteRepository(url string, debug, allowInsecure bool) *resolver.Repository {
	repo := resolver.NewRemoteRepository(url, debug, allowInsecure)
	return &repo
}

func getSQLiteRepository(dsn string) (resolver.Repository, error) {
	repo, err := resolver.NewSqliteRepository(dsn)
	return repo, err
}

func init() {
	Set("resolver", setupResolverService)
}

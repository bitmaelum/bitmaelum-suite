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
	"errors"
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
)

// We can have multiple resolvers to resolve a single address. We could resolve locally, remotely through
// resolver-services, or through DHT. We chain them all together with the ChainRepository

const defaultRemoteURL = "https://resolver.bitmaelum.com"

var (
	resolveOnce    sync.Once
	resolveService *resolver.Service
)

func setupResolverService() (interface{}, error) {
	var repo resolver.Repository
	var err error

	resolveOnce.Do(func() {
		switch config.Server.DefaultResolver {
		default:
			fallthrough
		case "remote":
			repo, err = getRemoteRepository()
		case "sqlite":
			repo, err = getSQLiteRepository()
		case "chain":
			repo, err = getChainRepository()
		}

		if err != nil {
			return
		}

		resolveService = resolver.KeyRetrievalService(repo)
	})

	// No correct resolver found, default to remote resolver
	if resolveService == nil {
		repo, err = getRemoteRepository()
		if err != nil {
			return nil, err
		}

		resolveService = resolver.KeyRetrievalService(repo)
	}

	return resolveService, nil
}

func getRemoteRepository() (resolver.Repository, error) {
	var (
		url           string
		debug         bool
		allowInsecure bool
	)

	if config.IsLoaded("client") {
		url = config.Client.Resolvers.Remote.URL
		debug = config.Client.Server.DebugHTTP
		allowInsecure = config.Client.Resolvers.Remote.AllowInsecure
	}
	if config.IsLoaded("server") {
		url = config.Server.Resolvers.Remote.URL
		debug = false
		allowInsecure = config.Server.Resolvers.Remote.AllowInsecure
	}
	if config.IsLoaded("bridge") {
		url = config.Bridge.Resolvers.Remote.URL
		debug = false
		allowInsecure = config.Bridge.Resolvers.Remote.AllowInsecure
	}

	// Set default URL if none is given
	if url == "" {
		url = defaultRemoteURL
	}

	return resolver.NewRemoteRepository(url, debug, allowInsecure), nil
}

func getSQLiteRepository() (resolver.Repository, error) {
	var path string

	if config.IsLoaded("client") {
		path = config.Client.Resolvers.Sqlite.Path
	}
	if config.IsLoaded("server") {
		path = config.Server.Resolvers.Sqlite.Path
	}
	if config.IsLoaded("bridge") {
		path = config.Bridge.Resolvers.Sqlite.Path
	}

	return resolver.NewSqliteRepository(path)
}

func getChainRepository() (resolver.Repository, error) {
	repo := resolver.NewChainRepository()

	var resolvers []string
	if config.IsLoaded("client") {
		resolvers = config.Client.Resolvers.Chain
	}
	if config.IsLoaded("server") {
		resolvers = config.Server.Resolvers.Chain
	}
	if config.IsLoaded("bridge") {
		resolvers = config.Bridge.Resolvers.Chain
	}

	var (
		chainedRepo resolver.Repository
		err         error
	)

	idx := 0
	for _, resolverName := range resolvers {
		switch resolverName {
		default:
			fallthrough
		case "remote":
			chainedRepo, err = getRemoteRepository()
		case "sqlite":
			chainedRepo, err = getSQLiteRepository()
		}

		if err != nil {
			return nil, err
		}

		_ = repo.(*resolver.ChainRepository).Add(chainedRepo)
		idx++
	}

	if idx == 0 {
		return nil, errors.New("the chain repo is empty")
	}

	return repo, nil
}

func init() {
	Instance.SetShared("resolver", setupResolverService)
}

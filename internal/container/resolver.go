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
	"errors"
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
	var repo resolver.Repository
	var err error

	resolveOnce.Do(func() {
		switch config.Server.DefaultResolver {
		default:
			fallthrough
		case "remote":
			repo, err = getRepo("remote")
		case "sqlite":
			repo, err = getRepo("sqlite")
		case "chain":
			repo, err = getRepo("chain")
		}

		if err != nil {
			return
		}

		resolveService = resolver.KeyRetrievalService(repo)
	})

	// No correct resolver found, default to remote resolver
	if resolveService == nil {
		repo, err = getRepo("remote")
		if err != nil {
			return nil, err
		}

		resolveService = resolver.KeyRetrievalService(repo)
	}

	return resolveService, nil
}

func getRepo(repoName string) (resolver.Repository, error) {
	var (
		repo resolver.Repository
		err error
	)

	switch (repoName) {
	case "remote":
		var (
			url string
			debug bool
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

		repo, err = getRemoteRepository(url, debug, allowInsecure)

	case "sqlite":
		if config.IsLoaded("client") {
			repo, err = getSQLiteRepository(config.Client.Resolvers.Sqlite.Path)
		}
		if config.IsLoaded("server") {
			repo, err = getSQLiteRepository(config.Server.Resolvers.Sqlite.Path)
		}
		if config.IsLoaded("bridge") {
			repo, err = getSQLiteRepository(config.Bridge.Resolvers.Sqlite.Path)
		}

	case "chain":
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

		idx := 0
		for _, resolverName := range resolvers {
			chainedRepo, err := getRepo(resolverName)
			if err != nil {
				return nil, err
			}

			_ = repo.(*resolver.ChainRepository).Add(chainedRepo)
			idx++
		}

		if idx == 0 {
			return nil, errors.New("the chain repo is empty")
		}
	}

	if err != nil {
		return nil, err
	}

	if repo == nil {
		return nil, errors.New("resolver not correctly configured")
	}

	return repo, nil
}

func getRemoteRepository(url string, debug, allowInsecure bool) (resolver.Repository, error) {
	repo := resolver.NewRemoteRepository(url, debug, allowInsecure)
	return repo, nil
}

func getSQLiteRepository(dsn string) (resolver.Repository, error) {
	repo, err := resolver.NewSqliteRepository(dsn)
	return repo, err
}

func init() {
	Instance.SetShared("resolver", setupResolverService)
}

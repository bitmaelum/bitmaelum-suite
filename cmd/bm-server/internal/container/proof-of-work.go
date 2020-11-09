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
	"os"
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/pow"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/go-redis/redis/v8"
)

var (
	proofOnce    sync.Once
	proofService pow.Storable
)

func setupProofService() (interface{}, error) {
	proofOnce.Do(func() {
		//If redis.host is set on the config file it will use redis instead of bolt
		if config.Server.Redis.Host != "" {
			opts := redis.Options{
				Addr: config.Server.Redis.Host,
				DB:   config.Server.Redis.Db,
			}

			proofService = pow.NewRedis(&opts)
		} else {
			if config.Server.Bolt.DatabasePath == "" {
				config.Server.Bolt.DatabasePath = os.TempDir()
			}

			proofService = pow.NewBolt(&config.Server.Bolt.DatabasePath)
		}

	})

	return proofService, nil
}

func init() {
	Instance.SetShared("proof-of-work", setupProofService)
}

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
	"os"
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"

	"github.com/go-redis/redis/v7"
)

var (
	webhookOnce       sync.Once
	webhookRepository *webhook.Repository
)

func setupWebhookRepo() (interface{}, error) {
	webhookOnce.Do(func() {
		// If redis.host is set on the config file it will use redis instead of bolt
		if config.Server.Redis.Host != "" {
			opts := redis.Options{
				Addr: config.Server.Redis.Host,
				DB:   config.Server.Redis.Db,
			}

			webhookRepository = webhook.NewRedisRepository(&opts)
			return
		}

		// If redis is not set then it will use BoltDB as default
		if config.Server.Bolt.DatabasePath == "" {
			config.Server.Bolt.DatabasePath = os.TempDir()
		}

		webhookRepository = webhook.NewBoltRepository(config.Server.Bolt.DatabasePath)
	})

	return *webhookRepository, nil
}

func init() {
	Instance.SetShared("webhook", setupWebhookRepo)
}

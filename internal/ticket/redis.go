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

package ticket

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type redisRepo struct {
	client *redis.Client
}

// NewRedisRepository initializes a new repository
func NewRedisRepository(opts *redis.Options) Repository {
	return &redisRepo{
		client: redis.NewClient(opts),
	}
}

// Fetch a ticket from the repository, or err
func (r redisRepo) Fetch(ticketID string) (*Ticket, error) {
	logrus.Trace("Trying to fetch ticket from REDIS: ", ticketID)
	data, err := r.client.Get(r.client.Context(), createTicketKey(ticketID)).Result()
	if data == "" || err != nil {
		logrus.Trace("ticket not found in REDIS: ", data, err)
		return nil, errors.New("ticket not found")
	}

	ticket := &Ticket{}
	err = json.Unmarshal([]byte(data), &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// Store the given ticket in the repository
func (r redisRepo) Store(ticket *Ticket) error {
	logrus.Trace("Storing ticket in REDIS: ", ticket)
	_, err := r.client.Set(r.client.Context(), createTicketKey(ticket.ID), ticket, 30*time.Minute).Result()

	return err
}

// Remove the given ticket from the repository
func (r redisRepo) Remove(ticketID string) {
	_ = r.client.Del(r.client.Context(), createTicketKey(ticketID))
}

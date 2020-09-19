package ticket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

type redisRepo struct {
	client *redis.Client
}

// NewRepository initializes a new repository
func NewRepository(opts *redis.Options) Repository {
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

// createTicketKey creates a key based on the given ID. This is needed otherwise we might send any data as ticket-id
// to redis in order to extract other kind of data (and you don't want that).
func createTicketKey(id string) string {
	return fmt.Sprintf("ticket-%s", id)
}

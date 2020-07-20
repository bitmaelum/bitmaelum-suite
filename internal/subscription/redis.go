package subscription

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis/v8"
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

func (r redisRepo) Has(sub *Subscription) bool {
	i, err := r.client.Exists(r.client.Context(), createKey(sub)).Result()

	return err == nil && i > 0
}

func (r redisRepo) Store(sub *Subscription) error {
	_, err := r.client.Set(r.client.Context(), createKey(sub), sub, 0).Result()

	return err
}

func (r redisRepo) Remove(sub *Subscription) error {
	_, err := r.client.Del(r.client.Context(), createKey(sub)).Result()

	return err
}

// Generate a key that can be used for reading / writing the subscription info
func createKey(sub *Subscription) string {
	data := fmt.Sprintf("%s-%s-%s", sub.From.String(), sub.To.String(), sub.SubscriptionID)

	h := sha256.New()
	return "subscription-" + hex.EncodeToString(h.Sum([]byte(data)))
}

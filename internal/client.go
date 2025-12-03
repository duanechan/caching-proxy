package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type client struct {
	timeout time.Duration
	rdb     *redis.Client
}

func newClient(timeout time.Duration) *client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &client{
		timeout: timeout,
		rdb:     rdb,
	}
}

func (c client) Add(key string, entry *CacheEntry) error {
	value, err := json.Marshal(*entry)
	if err != nil {
		return err
	}

	if err := c.rdb.Set(context.Background(), key, value, c.timeout).Err(); err != nil {
		return err
	}

	return nil
}

func (c client) Get(key string) (*CacheEntry, error) {
	value, err := c.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(value), &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

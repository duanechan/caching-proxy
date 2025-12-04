package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type client struct {
	timeout time.Duration
	Rdb     *redis.Client
}

func NewClient(timeout time.Duration) *client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pubsub := rdb.Subscribe(context.Background(), "proxy_control")
	go func() {
		for msg := range pubsub.Channel() {
			switch msg.Payload {
			case "clear-cache":
				ServerLog(Bold + "--- CACHE CLEARED" + Reset)
			case "get-cache":
				ServerLog(Bold + "--- CACHE HIT" + Reset)
			case "add-cache":
				ServerLog(Bold + "--- CACHE MISS" + Reset)
			default:
				WarnLog("----- UNKNOWN EVENT:", msg.Payload)
			}
		}
	}()

	return &client{
		timeout: timeout,
		Rdb:     rdb,
	}
}

func (c client) Add(key string, entry *CacheEntry) error {
	value, err := json.Marshal(*entry)
	if err != nil {
		return err
	}

	if err := c.Rdb.Set(context.Background(), key, value, 0).Err(); err != nil {
		return err
	}

	c.Rdb.Publish(context.Background(), "proxy_control", "add-cache")

	return nil
}

func (c client) Get(key string) (*CacheEntry, error) {
	value, err := c.Rdb.Get(context.Background(), key).Result()
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

	c.Rdb.Publish(context.Background(), "proxy_control", "get-cache")

	return &entry, nil
}

func (c client) FlushCache() error {
	if err := c.Rdb.FlushAll(context.Background()).Err(); err != nil {
		return err
	}

	fmt.Println("Cache cleared.")
	c.Rdb.Publish(context.Background(), "proxy_control", "clear-cache")

	return nil
}

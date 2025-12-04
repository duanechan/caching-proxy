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
				ServerLog(Bold+"-----", Green+"CACHE CLEARED"+Reset, Bold+"-----"+Reset)
			case "get-cache":
				ServerLog(Bold+"-----", Blue+"CACHE HIT"+Reset, Bold+"-----"+Reset)
			case "add-cache":
				ServerLog(Bold+"-----", Red+"CACHE MISS"+Reset, Bold+"-----"+Reset)
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

func (c client) add(key string, entry *CacheEntry) error {
	value, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if err := c.Rdb.Set(context.Background(), key, value, 0).Err(); err != nil {
		return err
	}

	c.Rdb.Publish(context.Background(), "proxy_control", "add-cache")

	return nil
}

func (c client) get(key string) (*CacheEntry, error) {
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
	keys, err := c.Rdb.DBSize(context.Background()).Result()
	if err != nil {
		return err
	}

	if keys == 0 {
		fmt.Println(Bold + Yellow + "No cache to flush." + Reset)
		return nil
	}

	if err := c.Rdb.FlushAll(context.Background()).Err(); err != nil {
		return err
	}

	fmt.Println(Bold+Green+"Cache cleared!"+Reset, keys, "keys removed.")
	c.Rdb.Publish(context.Background(), "proxy_control", "clear-cache")

	return nil
}

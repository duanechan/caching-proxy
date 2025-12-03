package proxy

import (
	"io"
	"net/http"
	"time"
)

type CacheEntry struct {
	Headers   http.Header `json:"headers"`
	Body      []byte      `json:"body"`
	CreatedAt time.Time   `json:"createdAt"`
}

func newCacheEntry(r *http.Response) (*CacheEntry, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return &CacheEntry{
		Headers:   r.Header,
		Body:      body,
		CreatedAt: time.Now(),
	}, nil
}

package proxy

import (
	"io"
	"net/http"
	"time"
)

type CacheEntry struct {
	StatusCode int         `json:"statusCode"`
	Body       []byte      `json:"body"`
	Headers    http.Header `json:"headers"`
	CreatedAt  time.Time   `json:"createdAt"`
}

func newCacheEntry(r *http.Response) (*CacheEntry, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return &CacheEntry{
		StatusCode: r.StatusCode,
		Headers:    r.Header.Clone(),
		Body:       body,
		CreatedAt:  time.Now(),
	}, nil
}

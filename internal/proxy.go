package proxy

import (
	"fmt"
	"net/http"
	"time"
)

type proxy struct {
	port   int
	origin string
	client client
}

func NewProxy(origin string, port int) (http.Handler, error) {
	url, err := normalizeURL(origin)
	if err != nil {
		return nil, err
	}

	return proxy{
		port:   port,
		origin: url,
		client: *newClient(5 * time.Second),
	}, nil
}

func (p proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entry, err := p.client.Get(r.URL.Path)
	if err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}

	if entry != nil {
		w.Header().Set("Content-Type", entry.Headers["Content-Type"][0])
		w.Header().Set("X-Cache", "HIT")
		fmt.Fprint(w, string(entry.Body))
		return
	}

	fullURL := p.origin + r.URL.Path
	res, err := http.Get(fullURL)
	if err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}

	entry, err = newCacheEntry(res)
	if err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}

	w.Header().Set("Content-Type", entry.Headers["Content-Type"][0])
	w.Header().Set("X-Cache", "MISS")
	fmt.Fprint(w, string(entry.Body))
}

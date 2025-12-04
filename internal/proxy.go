package proxy

import (
	"net/http"
	"time"
)

type proxy struct {
	port   int
	origin string
	client client
}

func NewProxy(client *client, origin string, port int) (http.Handler, error) {
	url, err := normalizeURL(origin)
	if err != nil {
		return nil, err
	}

	return proxy{
		port:   port,
		origin: url,
		client: *client,
	}, nil
}

func (p proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.URL.RequestURI()
	fullURL := p.origin + cacheKey

	start := time.Now()
	entry, err := p.client.Get(cacheKey)
	end := time.Since(start)
	if err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}

	if entry != nil {
		ServerLog(formatRequestLog(fullURL, end.Milliseconds()))
		w.Header().Set("X-Cache", "HIT")
		setHeaders(w, entry)
		return
	}

	start = time.Now()
	res, err := http.Get(fullURL)
	end = time.Since(start)
	if err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}
	defer res.Body.Close()

	entry, err = newCacheEntry(res)
	if err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}

	if err := p.client.Add(cacheKey, entry); err != nil {
		errorResponse(w, 400, "Bad Request")
		return
	}

	ServerLog(formatRequestLog(fullURL, end.Milliseconds()))
	w.Header().Set("X-Cache", "MISS")
	setHeaders(w, entry)
}

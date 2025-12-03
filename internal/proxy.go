package proxy

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

type proxy struct {
	port   int
	origin string
	rdb    *redis.Client
}

func NewProxy(origin string, port int) (http.Handler, error) {
	url, err := normalizeURL(origin)
	if err != nil {
		return nil, err
	}

	return proxy{
		port:   port,
		origin: url,
	}, nil
}

func (p proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

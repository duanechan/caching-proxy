package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("missing url argument")
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("invalid URL")
	}

	if parsedURL.Host == "" {
		return "", errors.New("missing host")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

	normalized := parsedURL.Scheme + "://" + strings.ToLower(parsedURL.Host)

	return normalized, nil
}

func setHeaders(w http.ResponseWriter, entry *CacheEntry) {
	for k, values := range entry.Headers {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}

	if entry.StatusCode > 0 {
		w.WriteHeader(entry.StatusCode)
	}

	w.Write(entry.Body)
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))
	w.WriteHeader(code)
}

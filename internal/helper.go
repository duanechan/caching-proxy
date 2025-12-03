package proxy

import (
	"errors"
	"fmt"
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

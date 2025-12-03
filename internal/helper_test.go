package proxy

import (
	"errors"
	"fmt"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{
			input:    "https://dummyjson.com/",
			expected: "https://dummyjson.com",
			err:      nil,
		},
		{
			input:    "dummyjson.com/",
			expected: "http://dummyjson.com",
			err:      nil,
		},
		{
			input:    "http://example.com/path",
			expected: "http://example.com",
			err:      nil,
		},
		{
			input:    "example.com/path",
			expected: "http://example.com",
			err:      nil,
		},
		{
			input:    "example.com:8080",
			expected: "http://example.com:8080",
			err:      nil,
		},
		{
			input:    "https://example.com:443/path?query=1",
			expected: "https://example.com:443",
			err:      nil,
		},
		{
			input:    "example.com?foo=bar",
			expected: "http://example.com",
			err:      nil,
		},
		{
			input:    "ftp://example.com",
			expected: "",
			err:      fmt.Errorf("unsupported scheme: ftp"),
		},
		{
			input:    "://example.com",
			expected: "",
			err:      fmt.Errorf("invalid URL"),
		},
		{
			input:    "http://example.com/",
			expected: "http://example.com",
			err:      nil,
		},
		{
			input:    "example.com/",
			expected: "http://example.com",
			err:      nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("basic test #%d", i+1), func(t *testing.T) {
			actual, err := normalizeURL(test.input)
			if test.err != nil {
				if !errors.Is(err, test.err) && err.Error() != test.err.Error() {
					t.Errorf("expected err to be %v, got %v", test.err, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if actual != test.expected {
				t.Errorf("expected url to be %s, got %s", test.expected, actual)
			}
		})
	}
}

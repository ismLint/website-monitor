package storage

import (
	"fmt"
	"net/url"
	"strings"
)

func ValidateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return fmt.Errorf("%w: URL cannot be empty", ErrInvalidURL)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("%w: Failed to parse URL: %v", ErrInvalidURL, err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: Can't use such things site", ErrInvalidURL)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("%w: URL host cannot be empty", ErrInvalidURL)
	}

	return nil
}

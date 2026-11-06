package storage

import (
	"errors"
	"net/url"
	"strings"
)

func ValidateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return errors.New("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.New("Error parse URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("Can't use such things site")
	}

	if parsedURL.Host == "" {
		return errors.New("URL host cannot be empty")
	}

	return nil
}

package checker

import (
	"net/http"
	"time"
)

type CheckResult struct {
	URL            string
	HTTPStatusCode int
	Error          error
	ElapsedTime    time.Duration
	Available      bool
}

func CheckWebsite(url string) CheckResult {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	res, err := client.Get(url)
	if err != nil {
		return CheckResult{
			URL:         url,
			Error:       err,
			ElapsedTime: time.Since(start),
			Available:   false,
		}
	}
	defer res.Body.Close()

	return CheckResult{
		URL:            url,
		HTTPStatusCode: res.StatusCode,
		Error:          nil,
		ElapsedTime:    time.Since(start),
		Available:      true,
	}
}

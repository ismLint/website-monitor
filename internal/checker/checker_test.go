package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckWebsite(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
	}{
		{
			name:          "Valid website",
			statusCode:    http.StatusOK,
			expectedError: false,
		},
		{
			name:          "Server error",
			statusCode:    http.StatusInternalServerError,
			expectedError: false,
		},
		{
			name:          "Not found",
			statusCode:    http.StatusNotFound,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
				},
			))

			defer server.Close()

			result := CheckWebsite(server.URL)

			if result.URL != server.URL {
				t.Errorf("Expected URL %s, got %s",
					server.URL, result.URL)
			}

			if (result.Error != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v",
					tt.expectedError, result.Error)
			}

			if result.ElapsedTime <= 0 {
				t.Errorf("Expected positive elapsed time, got %v",
					result.ElapsedTime)
			}

			if result.HTTPStatusCode != tt.statusCode {
				t.Errorf("Expected status %d, got %d",
					tt.statusCode, result.HTTPStatusCode)
			}
		})
	}
}

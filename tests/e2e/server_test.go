package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"website-monitor/internal/app"
	"website-monitor/internal/storage"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func TestHealthEndpoint(t *testing.T) {
	store := storage.NewStorage()
	router := app.NewRouter(store)
	var result HealthResponse

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d status code, got %d", http.StatusOK, resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Parse JSON failed: %v", err)
	}

	if result.Service != "website-monitor" || result.Status != "ok" {
		t.Fatalf("Expected service 'website-monitor', got service %v. Expected status 'ok', got %v",
			result.Service, result.Status)
	}

}

func TestGetSitesEndpoint(t *testing.T) {
	store := storage.NewStorage()
	router := app.NewRouter(store)
	var sites []storage.Site

	if _, err := store.Add(storage.Site{URL: "https://example.com"}); err != nil {
		t.Fatalf("Failed to add site: %v", err)
	}
	if _, err := store.Add(storage.Site{URL: "https://google.com"}); err != nil {
		t.Fatalf("Error to add site: %v", err)
	}

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/sites")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&sites); err != nil {
		t.Fatalf("Parse JSON failed: %v", err)
	}

	if len(sites) != 2 {
		t.Fatalf("Expected 2 sites, got: %d", len(sites))
	}

	expected := map[string]bool{
		"https://example.com": false,
		"https://google.com":  false,
	}

	for _, site := range sites {
		if _, ok := expected[site.URL]; !ok {
			t.Fatalf("Unexpected site URL returned: %s", site.URL)
		}
		expected[site.URL] = true
	}
	for url, found := range expected {
		if !found {
			t.Fatalf("Expected site %s to be returned", url)
		}
	}
}

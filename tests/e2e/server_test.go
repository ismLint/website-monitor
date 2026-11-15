package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestGetSiteByIDEndpoint(t *testing.T) {
	store := storage.NewStorage()
	router := app.NewRouter(store)

	var site storage.Site

	createdSite1, err := store.Add(storage.Site{URL: "https://example.com"})
	if err != nil {
		t.Fatalf("Failed to add site: %v", err)
	}

	if _, err := store.Add(storage.Site{URL: "https://google.com"}); err != nil {
		t.Fatalf("Failed to add site: %v", err)
	}

	server := httptest.NewServer(router)
	defer server.Close()

	req := fmt.Sprintf("%s/sites/%d", server.URL, createdSite1.ID)
	resp, err := http.Get(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&site); err != nil {
		t.Fatalf("Parse JSON failed: %v", err)
	}

	if createdSite1.URL != site.URL {
		t.Fatalf("Expected %s url, got %v ", createdSite1.URL, site.URL)
	}

	if createdSite1.ID != site.ID {
		t.Fatalf("Expected %d id, got %d", createdSite1.ID, site.ID)
	}
}

func TestGetByIDEndpointTable(t *testing.T) {
	store := storage.NewStorage()
	router := app.NewRouter(store)
	var site storage.Site

	createdSite, err := store.Add(storage.Site{URL: "http://example.com"})
	if err != nil {
		t.Fatalf("Failed to add site: %v", err)
	}

	server := httptest.NewServer(router)
	defer server.Close()

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedID     int
	}{
		{
			name:           "Site Not Found",
			id:             "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid ID",
			id:             "qwe",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid ID",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedID:     createdSite.ID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := fmt.Sprintf("%s/sites/%s", server.URL, tt.id)
			resp, err := http.Get(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("Expected %d status code, got %d. Test name: %s",
					tt.expectedStatus, resp.StatusCode, tt.name)
			}

			if resp.StatusCode == http.StatusOK {
				if err = json.NewDecoder(resp.Body).Decode(&site); err != nil {
					t.Fatalf("Failed to decode JSON: %v", err)
				}

				if tt.expectedID != site.ID {
					t.Fatalf("Expected %d site id, got %d", tt.expectedID, site.ID)
				}
			}
		})
	}
}

func TestPostSitesEncpoint(t *testing.T) {
	var createdSite storage.Site
	requestBody := struct {
		URL string `json:"url"`
	}{
		URL: "https://example.com",
	}
	store := storage.NewStorage()
	router := app.NewRouter(store)
	server := httptest.NewServer(router)
	defer server.Close()

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed serialize struct: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		server.URL+"/sites",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Request Failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Response send failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d",
			http.StatusCreated, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&createdSite)
	if err != nil {
		t.Fatalf("Failed decode JSON: %v", err)
	}

	if createdSite.ID < 1 {
		t.Fatal("Site ID < 1")
	}

	if requestBody.URL != createdSite.URL {
		t.Fatalf("Expected url %s, got %s",
			req.URL.String(), createdSite.URL)
	}
}

// func TestPatchSitesEndpoint(t *testing.T) {
// 	requestBody := struct {
// 		URL string `json:"url"`
// 		ID  int    `json:"id"`
// 	}{
// 		ID:  1,
// 		URL: "https://google.com",
// 	}
// 	store := storage.NewStorage()
// 	router := app.NewRouter(store)

// 	server := httptest.NewServer(router)
// 	defer server.Close()

// 	createdSite, err := store.Add(storage.Site{URL: "https://example.com"})
// 	if err != nil {
// 		t.Fatalf("Failed add site: %v", err)
// 	}

// 	jsonData, err := json.Marshal(requestBody)
// 	if err != nil {
// 		t.Fatalf("Failed serialize struct: %v", err)
// 	}

// 	request := fmt.Sprintf("%s/sites/%d", server.URL, requestBody.ID)

// 	req, err := http.NewRequest(
// 		http.MethodPatch,
// 		request,
// 		bytes.NewBuffer(jsonData),
// 	)

// }

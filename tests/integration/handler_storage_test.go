package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"website-monitor/internal/handler"
	"website-monitor/internal/storage"
)

func TestPostSiteIntegration(t *testing.T) {
	store := storage.NewStorage()
	h := handler.NewHandler(store)

	requestBody := struct {
		URL string `json:"url"`
	}{
		URL: "https://example.com",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/sites",
		bytes.NewBuffer(jsonData),
	)

	req.Header.Set("Content-Type", "application/json")
	h.PostSites(rec, req)
	defer req.Body.Close()

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d", rec.Code)
	}

	sites := store.GetAll()

	if len(sites) != 1 {
		t.Fatalf("Expected 1 site in storage, got %d", len(sites))
	}

	if sites[0].URL != "https://example.com" {
		t.Fatalf("Expected stored URL to be https://example.com, got %s", sites[0].URL)
	}

	if sites[0].ID <= 0 {
		t.Fatalf("Expected positive site ID, got %d", sites[0].ID)
	}
}

func TestPatchSiteIntegration(t *testing.T) {
	store := storage.NewStorage()
	h := handler.NewHandler(store)

	createdSite, err := store.Add(storage.Site{URL: "http://example.com"})
	if err != nil {
		t.Fatalf("Failed to add site: %v", err)
	}

	requestBody := struct {
		URL string `json:"url"`
	}{
		URL: "https://google.com",
	}

	requestURL := fmt.Sprintf("/sites/%d", createdSite.ID)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPatch,
		requestURL,
		bytes.NewBuffer(jsonData),
	)
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")

	h.PatchSites(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected %d status code, got %d", http.StatusOK, rec.Code)
	}

	site, err := store.GetByID(createdSite.ID)
	if err != nil {
		t.Fatalf("Failed to get site by id: %v", err)
	}

	if site.URL != requestBody.URL {
		t.Fatalf("Expected %s, got %s", requestBody.URL, site.URL)
	}

	if site.Availability {
		t.Fatalf("Expected availabity false, got true")
	}

	if !site.LastCheck.IsZero() {
		t.Fatalf("Expected zero last check")
	}

	if site.ResponseTime != 0 {
		t.Fatalf("Expected zero response time")
	}

	if site.ID != createdSite.ID {
		t.Fatalf("Expected ID %d, got %d", createdSite.ID, site.ID)
	}
}

func TestDeleteSiteIntegration(t *testing.T) {
	store := storage.NewStorage()
	h := handler.NewHandler(store)

	createdSite, err := store.Add(storage.Site{URL: "https://example.com"})
	if err != nil {
		t.Fatalf("Failed to add site: %v", err)
	}

	requestURL := fmt.Sprintf("/sites/%d", createdSite.ID)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodDelete,
		requestURL,
		nil,
	)

	h.DeleteSite(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("Expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	_, err = store.GetByID(createdSite.ID)
	if !errors.Is(err, storage.ErrSiteNotFound) {
		t.Fatalf("Expected Not found, got %v", err)
	}

}

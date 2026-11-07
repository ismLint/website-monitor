package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"website-monitor/internal/handler"
	"website-monitor/internal/storage"
)

type Server struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func main() {
	mux := http.NewServeMux()

	store := storage.NewStorage()
	h := handler.NewHandler(store)

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /sites", h.GetAllSites)
	mux.HandleFunc("GET /sites/{id}", h.GetSites)
	mux.HandleFunc("POST /sites", h.PostSites)
	mux.HandleFunc("PATCH /sites/{id}", h.PatchSites)
	mux.HandleFunc("DELETE /sites/{id}", h.DeleteSite)

	store.Add(storage.Site{URL: "https://example.com"})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Server{
		Status:  "ok",
		Service: "website-monitor",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

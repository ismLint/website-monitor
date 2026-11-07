package app

import (
	"net/http"

	"website-monitor/internal/handler"
	"website-monitor/internal/storage"
)

func NewRouter(store *storage.Storage) http.Handler {
	mux := http.NewServeMux()
	h := handler.NewHandler(store)

	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("GET /sites", h.GetAllSites)
	mux.HandleFunc("GET /sites/{id}", h.GetSites)
	mux.HandleFunc("POST /sites", h.PostSites)
	mux.HandleFunc("PATCH /sites/{id}", h.PatchSites)
	mux.HandleFunc("DELETE /sites/{id}", h.DeleteSite)

	return mux
}

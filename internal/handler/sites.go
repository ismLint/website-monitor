package handler

import (
	"encoding/json"
	"net/http"
	"website-monitor/internal/storage"
)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(s *storage.Storage) *Handler {
	return &Handler{
		storage: s,
	}
}

func (h *Handler) GetAllSites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sites := h.storage.GetAll()

	if err := json.NewEncoder(w).Encode(sites); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

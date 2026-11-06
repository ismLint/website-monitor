package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

func (h *Handler) GetSites(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	site, err := h.storage.GetByID(id)
	if errors.Is(err, storage.ErrSiteNotFound) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(site); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

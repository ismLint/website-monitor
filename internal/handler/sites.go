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

var request struct {
	URL string `json:"url"`
}

func NewHandler(s *storage.Storage) *Handler {
	return &Handler{
		storage: s,
	}
}

type Server struct {
	Status  string `json:"status"`
	Service string `json:"service"`
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

func (h *Handler) PostSites(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 mib
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	site, err := h.storage.Add(storage.Site{URL: request.URL})
	if errors.Is(err, storage.ErrInvalidURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(site); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PatchSites(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	site, err := h.storage.Update(id, request.URL)
	if errors.Is(err, storage.ErrInvalidURL) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if errors.Is(err, storage.ErrSiteNotFound) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(site); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.storage.Delete(id)
	if errors.Is(err, storage.ErrSiteNotFound) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
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

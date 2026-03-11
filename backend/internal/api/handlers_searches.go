package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"search-mm2/backend/internal/database"
	"search-mm2/backend/internal/models"
	"search-mm2/backend/internal/scraper"
)

type SearchHandlers struct {
	db      *database.Queries
	scraper *scraper.Service
}

func NewSearchHandlers(db *database.Queries, s *scraper.Service) *SearchHandlers {
	return &SearchHandlers{db: db, scraper: s}
}

func (h *SearchHandlers) List(w http.ResponseWriter, r *http.Request) {
	searches, err := h.db.ListSearches(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list searches"}`, http.StatusInternalServerError)
		return
	}
	if searches == nil {
		searches = []models.Search{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searches)
}

func (h *SearchHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var s models.Search
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(s.URL) == "" {
		http.Error(w, `{"error":"url is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.CreateSearch(r.Context(), &s); err != nil {
		http.Error(w, `{"error":"failed to create search"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *SearchHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	s, err := h.db.GetSearch(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"search not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *SearchHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	var s models.Search
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	s.ID = id
	if strings.TrimSpace(s.URL) == "" {
		http.Error(w, `{"error":"url is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateSearch(r.Context(), &s); err != nil {
		http.Error(w, `{"error":"failed to update search"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *SearchHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteSearch(r.Context(), id); err != nil {
		http.Error(w, `{"error":"failed to delete search"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SearchHandlers) Scrape(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	s, err := h.db.GetSearch(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"search not found"}`, http.StatusNotFound)
		return
	}

	go h.scraper.RunSearch(r.Context(), s)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"scrape started"}`))
}

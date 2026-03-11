package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"search-mm2/backend/internal/database"
	"search-mm2/backend/internal/models"
)

type PropertyHandlers struct {
	db *database.Queries
}

func NewPropertyHandlers(db *database.Queries) *PropertyHandlers {
	return &PropertyHandlers{db: db}
}

func (h *PropertyHandlers) ListBySearch(w http.ResponseWriter, r *http.Request) {
	searchID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid search id"}`, http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	props, err := h.db.ListPropertiesBySearch(r.Context(), searchID, limit, offset)
	if err != nil {
		http.Error(w, `{"error":"failed to list properties"}`, http.StatusInternalServerError)
		return
	}
	if props == nil {
		props = []models.Property{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(props)
}

func (h *PropertyHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	p, err := h.db.GetProperty(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"property not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

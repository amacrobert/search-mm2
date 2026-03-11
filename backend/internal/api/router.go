package api

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"search-mm2/backend/internal/config"
	"search-mm2/backend/internal/database"
	"search-mm2/backend/internal/scraper"
)

func NewRouter(cfg *config.Config, db *database.Queries, s *scraper.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	searchHandlers := NewSearchHandlers(db, s)
	propertyHandlers := NewPropertyHandlers(db)

	r.Post("/api/auth/login", HandleLogin(cfg))

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(cfg.JWTSecret))

		r.Get("/api/searches", searchHandlers.List)
		r.Post("/api/searches", searchHandlers.Create)
		r.Get("/api/searches/{id}", searchHandlers.Get)
		r.Put("/api/searches/{id}", searchHandlers.Update)
		r.Delete("/api/searches/{id}", searchHandlers.Delete)
		r.Post("/api/searches/{id}/scrape", searchHandlers.Scrape)
		r.Get("/api/searches/{id}/properties", propertyHandlers.ListBySearch)
		r.Get("/api/properties/{id}", propertyHandlers.Get)
	})

	return r
}

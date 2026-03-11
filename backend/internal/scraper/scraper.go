package scraper

import (
	"context"
	"log"

	"search-mm2/backend/internal/database"
	"search-mm2/backend/internal/models"
)

type Service struct {
	db *database.Queries
}

func NewService(db *database.Queries) *Service {
	return &Service{db: db}
}

func (s *Service) RunAll(ctx context.Context) {
	searches, err := s.db.GetActiveSearches(ctx)
	if err != nil {
		log.Printf("scraper: failed to get active searches: %v", err)
		return
	}

	log.Printf("scraper: running %d active searches", len(searches))
	for i := range searches {
		s.RunSearch(ctx, &searches[i])
	}
}

func (s *Service) RunSearch(ctx context.Context, search *models.Search) {
	log.Printf("scraper: scraping search %q (%s)", search.Name, search.ID)

	properties, err := ScrapeLoopNet(search)
	if err != nil {
		log.Printf("scraper: failed to scrape search %s: %v", search.ID, err)
		return
	}

	log.Printf("scraper: found %d properties for search %q", len(properties), search.Name)
	for i := range properties {
		properties[i].SearchID = search.ID
		if err := s.db.UpsertProperty(ctx, &properties[i]); err != nil {
			log.Printf("scraper: failed to upsert property %q: %v", properties[i].ExternalID, err)
		}
	}
}

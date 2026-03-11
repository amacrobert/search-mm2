package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"search-mm2/backend/internal/models"
)

type Queries struct {
	pool *pgxpool.Pool
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{pool: pool}
}

// Searches

func (q *Queries) CreateSearch(ctx context.Context, s *models.Search) error {
	return q.pool.QueryRow(ctx,
		`INSERT INTO searches (name, url, active)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		s.Name, s.URL, s.Active,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (q *Queries) GetSearch(ctx context.Context, id uuid.UUID) (*models.Search, error) {
	s := &models.Search{}
	err := q.pool.QueryRow(ctx,
		`SELECT id, name, url, active, created_at, updated_at
		 FROM searches WHERE id = $1`, id,
	).Scan(&s.ID, &s.Name, &s.URL, &s.Active, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (q *Queries) ListSearches(ctx context.Context) ([]models.Search, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT id, name, url, active, created_at, updated_at
		 FROM searches ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searches []models.Search
	for rows.Next() {
		var s models.Search
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.Active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		searches = append(searches, s)
	}
	return searches, rows.Err()
}

func (q *Queries) UpdateSearch(ctx context.Context, s *models.Search) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE searches SET name=$1, url=$2, active=$3, updated_at=now()
		 WHERE id=$4`,
		s.Name, s.URL, s.Active, s.ID,
	)
	return err
}

func (q *Queries) DeleteSearch(ctx context.Context, id uuid.UUID) error {
	_, err := q.pool.Exec(ctx, `DELETE FROM searches WHERE id = $1`, id)
	return err
}

func (q *Queries) GetActiveSearches(ctx context.Context) ([]models.Search, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT id, name, url, active, created_at, updated_at
		 FROM searches WHERE active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searches []models.Search
	for rows.Next() {
		var s models.Search
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.Active, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		searches = append(searches, s)
	}
	return searches, rows.Err()
}

// Properties

func (q *Queries) UpsertProperty(ctx context.Context, p *models.Property) error {
	return q.pool.QueryRow(ctx,
		`INSERT INTO properties (search_id, external_id, name, address, city, state, zip, property_type, price, size_sqft, description, url, image_url, listed_date, scraped_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		 ON CONFLICT (external_id, search_id) DO UPDATE SET
		   name=EXCLUDED.name, address=EXCLUDED.address, city=EXCLUDED.city, state=EXCLUDED.state, zip=EXCLUDED.zip,
		   property_type=EXCLUDED.property_type, price=EXCLUDED.price, size_sqft=EXCLUDED.size_sqft,
		   description=EXCLUDED.description, url=EXCLUDED.url, image_url=EXCLUDED.image_url,
		   listed_date=EXCLUDED.listed_date, scraped_at=EXCLUDED.scraped_at
		 RETURNING id, created_at`,
		p.SearchID, p.ExternalID, p.Name, p.Address, p.City, p.State, p.Zip, p.PropertyType,
		p.Price, p.SizeSqFt, p.Description, p.URL, p.ImageURL, p.ListedDate, p.ScrapedAt,
	).Scan(&p.ID, &p.CreatedAt)
}

func (q *Queries) ListPropertiesBySearch(ctx context.Context, searchID uuid.UUID, limit, offset int) ([]models.Property, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := q.pool.Query(ctx,
		`SELECT id, search_id, external_id, name, address, city, state, zip, property_type, price, size_sqft, description, url, image_url, listed_date, scraped_at, created_at
		 FROM properties WHERE search_id = $1 ORDER BY scraped_at DESC LIMIT $2 OFFSET $3`,
		searchID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var props []models.Property
	for rows.Next() {
		var p models.Property
		if err := rows.Scan(&p.ID, &p.SearchID, &p.ExternalID, &p.Name, &p.Address, &p.City, &p.State, &p.Zip, &p.PropertyType, &p.Price, &p.SizeSqFt, &p.Description, &p.URL, &p.ImageURL, &p.ListedDate, &p.ScrapedAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		props = append(props, p)
	}
	return props, rows.Err()
}

func (q *Queries) GetProperty(ctx context.Context, id uuid.UUID) (*models.Property, error) {
	p := &models.Property{}
	err := q.pool.QueryRow(ctx,
		`SELECT id, search_id, external_id, name, address, city, state, zip, property_type, price, size_sqft, description, url, image_url, listed_date, scraped_at, created_at
		 FROM properties WHERE id = $1`, id,
	).Scan(&p.ID, &p.SearchID, &p.ExternalID, &p.Name, &p.Address, &p.City, &p.State, &p.Zip, &p.PropertyType, &p.Price, &p.SizeSqFt, &p.Description, &p.URL, &p.ImageURL, &p.ListedDate, &p.ScrapedAt, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get property: %w", err)
	}
	return p, nil
}

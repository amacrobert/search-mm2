package models

import (
	"time"

	"github.com/google/uuid"
)

type Search struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Property struct {
	ID           uuid.UUID  `json:"id"`
	SearchID     uuid.UUID  `json:"searchId"`
	ExternalID   string     `json:"externalId"`
	Name         string     `json:"name"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	Zip          string     `json:"zip"`
	PropertyType string     `json:"propertyType"`
	Price        *float64   `json:"price"`
	SizeSqFt     *int       `json:"sizeSqFt"`
	Description  string     `json:"description"`
	URL          string     `json:"url"`
	ImageURL     string     `json:"imageUrl"`
	ListedDate   *time.Time `json:"listedDate"`
	ScrapedAt    time.Time  `json:"scrapedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// Package models defines the data structures for the ecaytracker application.
package models

import "time"

// Listing represents a single car listing scraped from ecaytrade.com.
// Fields map 1-to-1 with the `listings` table in Supabase.
type Listing struct {
	ID           string     `json:"id,omitempty"`
	ExternalID   string     `json:"external_id"`
	URL          string     `json:"url"`
	Title        string     `json:"title"`
	Make         string     `json:"make,omitempty"`
	Model        string     `json:"model,omitempty"`
	Year         *int       `json:"year,omitempty"`
	Mileage      *int       `json:"mileage,omitempty"`
	Price        float64    `json:"price"`
	Currency     string     `json:"currency"`
	Condition    string     `json:"condition,omitempty"`
	Transmission string     `json:"transmission,omitempty"`
	FuelType     string     `json:"fuel_type,omitempty"`
	Color        string     `json:"color,omitempty"`
	Description  string     `json:"description,omitempty"`
	Images       []string   `json:"images,omitempty"`
	Location     string     `json:"location,omitempty"`
	SellerName   string     `json:"seller_name,omitempty"`
	IsActive     bool       `json:"is_active"`
	FirstSeen    *time.Time `json:"first_seen,omitempty"`
	LastSeen     *time.Time `json:"last_seen,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

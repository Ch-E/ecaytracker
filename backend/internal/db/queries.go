package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ecaytracker/backend/models"
)

// UpsertResult describes what happened when a listing was upserted.
type UpsertResult struct {
	Inserted     bool
	PriceChanged bool
}

// UpsertListing inserts a new listing or updates the existing one matched on
// external_id. If the price has changed, a price_history row is also inserted.
func UpsertListing(ctx context.Context, pool *pgxpool.Pool, l models.Listing) (UpsertResult, error) {
	var res UpsertResult

	// ── 1. Check whether the listing already exists and what its current price is ──
	var existingID string
	var existingPrice float64
	err := pool.QueryRow(ctx,
		`SELECT id, price FROM listings WHERE external_id = $1`,
		l.ExternalID,
	).Scan(&existingID, &existingPrice)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		res.Inserted = true
	case err != nil:
		return res, fmt.Errorf("check existing listing %s: %w", l.ExternalID, err)
	}

	// ── 2. Upsert ──
	var returnedID string
	err = pool.QueryRow(ctx, `
		INSERT INTO listings
			(external_id, url, title, make, model, year, mileage, price, currency,
			 images, location, is_active, last_seen)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,TRUE,NOW())
		ON CONFLICT (external_id) DO UPDATE SET
			url          = EXCLUDED.url,
			title        = EXCLUDED.title,
			make         = EXCLUDED.make,
			model        = EXCLUDED.model,
			year         = EXCLUDED.year,
			mileage      = EXCLUDED.mileage,
			price        = EXCLUDED.price,
			currency     = EXCLUDED.currency,
			images       = EXCLUDED.images,
			location     = EXCLUDED.location,
			is_active    = TRUE,
			last_seen    = NOW(),
			updated_at   = NOW()
		RETURNING id`,
		l.ExternalID, l.URL, l.Title, l.Make, l.Model,
		l.Year, l.Mileage, l.Price, l.Currency,
		l.Images, l.Location,
	).Scan(&returnedID)
	if err != nil {
		return res, fmt.Errorf("upsert listing %s: %w", l.ExternalID, err)
	}

	// ── 3. Record price history when the price has actually changed ──
	if !res.Inserted && existingPrice != l.Price && l.Price > 0 {
		res.PriceChanged = true
		_, err = pool.Exec(ctx,
			`INSERT INTO price_history (listing_id, price) VALUES ($1, $2)`,
			returnedID, l.Price,
		)
		if err != nil {
			return res, fmt.Errorf("insert price_history for %s: %w", l.ExternalID, err)
		}
	}

	return res, nil
}

// GetListings returns all active listings ordered newest-first.
func GetListings(ctx context.Context, pool *pgxpool.Pool) ([]models.Listing, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			id, external_id, url, title,
			make, model, year, mileage,
			price, currency, condition, transmission,
			fuel_type, color, description, images,
			location, seller_name, is_active,
			first_seen, last_seen, created_at, updated_at
		FROM listings
		WHERE is_active = TRUE
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	var listings []models.Listing
	for rows.Next() {
		var (
			l models.Listing
			// Nullable text columns — scan into *string, convert to string after.
			make_, model_, condition_, transmission_ *string
			fuelType_, color_, description_          *string
			location_, sellerName_                   *string
			// Nullable timestamptz columns.
			firstSeen_, lastSeen_, createdAt_, updatedAt_ *time.Time
		)

		err := rows.Scan(
			&l.ID, &l.ExternalID, &l.URL, &l.Title,
			&make_, &model_, &l.Year, &l.Mileage,
			&l.Price, &l.Currency, &condition_, &transmission_,
			&fuelType_, &color_, &description_, &l.Images,
			&location_, &sellerName_, &l.IsActive,
			&firstSeen_, &lastSeen_, &createdAt_, &updatedAt_,
		)
		if err != nil {
			return nil, fmt.Errorf("scan listing row: %w", err)
		}

		l.Make = strVal(make_)
		l.Model = strVal(model_)
		l.Condition = strVal(condition_)
		l.Transmission = strVal(transmission_)
		l.FuelType = strVal(fuelType_)
		l.Color = strVal(color_)
		l.Description = strVal(description_)
		l.Location = strVal(location_)
		l.SellerName = strVal(sellerName_)
		l.FirstSeen = firstSeen_
		l.LastSeen = lastSeen_
		l.CreatedAt = createdAt_
		l.UpdatedAt = updatedAt_

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return listings, nil
}

// strVal dereferences a *string, returning "" for nil pointers.
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

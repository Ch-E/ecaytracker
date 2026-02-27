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
			 images, location, condition, transmission, fuel_type, color,
			 body_type, drive, cylinders, steering, interior_color, doors, on_island,
			 is_active, last_seen)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,TRUE,NOW())
		ON CONFLICT (external_id) DO UPDATE SET
			url            = EXCLUDED.url,
			title          = EXCLUDED.title,
			make           = EXCLUDED.make,
			model          = EXCLUDED.model,
			year           = EXCLUDED.year,
			mileage        = EXCLUDED.mileage,
			price          = EXCLUDED.price,
			currency       = EXCLUDED.currency,
			images         = EXCLUDED.images,
			location       = EXCLUDED.location,
			condition      = EXCLUDED.condition,
			transmission   = EXCLUDED.transmission,
			fuel_type      = EXCLUDED.fuel_type,
			color          = EXCLUDED.color,
			body_type      = EXCLUDED.body_type,
			drive          = EXCLUDED.drive,
			cylinders      = EXCLUDED.cylinders,
			steering       = EXCLUDED.steering,
			interior_color = EXCLUDED.interior_color,
			doors          = EXCLUDED.doors,
			on_island      = EXCLUDED.on_island,
			is_active      = TRUE,
			last_seen      = NOW(),
			updated_at     = NOW()
		RETURNING id`,
		l.ExternalID, l.URL, l.Title, l.Make, l.Model,
		l.Year, l.Mileage, l.Price, l.Currency,
		l.Images, l.Location, l.Condition, l.Transmission, l.FuelType, l.Color,
		l.BodyType, l.Drive, l.Cylinders, l.Steering, l.InteriorColor, l.Doors, l.OnIsland,
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
			fuel_type, color, body_type, drive,
			cylinders, steering, interior_color, doors, on_island,
			description, images,
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
			// Nullable text columns.
			make_, model_, condition_, transmission_ *string
			fuelType_, color_, bodyType_, drive_     *string
			cylinders_, steering_, interiorColor_    *string
			doors_, description_, location_          *string
			sellerName_                              *string
			// Nullable timestamptz columns.
			firstSeen_, lastSeen_, createdAt_, updatedAt_ *time.Time
		)

		err := rows.Scan(
			&l.ID, &l.ExternalID, &l.URL, &l.Title,
			&make_, &model_, &l.Year, &l.Mileage,
			&l.Price, &l.Currency, &condition_, &transmission_,
			&fuelType_, &color_, &bodyType_, &drive_,
			&cylinders_, &steering_, &interiorColor_, &doors_, &l.OnIsland,
			&description_, &l.Images,
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
		l.BodyType = strVal(bodyType_)
		l.Drive = strVal(drive_)
		l.Cylinders = strVal(cylinders_)
		l.Steering = strVal(steering_)
		l.InteriorColor = strVal(interiorColor_)
		l.Doors = strVal(doors_)
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

// GetStats returns pre-computed dashboard statistics: total listing count, average
// price, median price, new-this-week count, average mileage, top 8 makes, body
// type distribution, and year distribution.
func GetStats(ctx context.Context, pool *pgxpool.Pool) (models.Stats, error) {
	var stats models.Stats

	// Single-row aggregates.
	err := pool.QueryRow(ctx, `
		SELECT
			COUNT(*)::int,
			COALESCE(AVG(price), 0),
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY price), 0),
			COUNT(*) FILTER (WHERE first_seen >= NOW() - INTERVAL '7 days')::int,
			COALESCE(AVG(mileage) FILTER (WHERE mileage IS NOT NULL), 0)
		FROM listings
		WHERE is_active = TRUE
	`).Scan(&stats.TotalListings, &stats.AvgPrice, &stats.MedianPrice, &stats.NewThisWeek, &stats.AvgMileage)
	if err != nil {
		return stats, fmt.Errorf("get stats aggregate: %w", err)
	}

	// Top 8 makes by listing count.
	brandRows, err := pool.Query(ctx, `
		SELECT make, COUNT(*)::int, COALESCE(AVG(price), 0)
		FROM listings
		WHERE is_active = TRUE AND make IS NOT NULL AND make != ''
		GROUP BY make
		ORDER BY COUNT(*) DESC
		LIMIT 8
	`)
	if err != nil {
		return stats, fmt.Errorf("get top brands: %w", err)
	}
	defer brandRows.Close()

	for brandRows.Next() {
		var b models.BrandStat
		if err := brandRows.Scan(&b.Name, &b.Count, &b.AvgPrice); err != nil {
			return stats, fmt.Errorf("scan brand row: %w", err)
		}
		stats.TopBrands = append(stats.TopBrands, b)
	}
	if err := brandRows.Err(); err != nil {
		return stats, fmt.Errorf("brand rows error: %w", err)
	}
	if stats.TopBrands == nil {
		stats.TopBrands = make([]models.BrandStat, 0)
	}

	// Body type distribution — null/empty values are grouped as "Other".
	btRows, err := pool.Query(ctx, `
		SELECT
			COALESCE(NULLIF(TRIM(body_type), ''), 'Other') AS bt,
			COUNT(*)::int,
			COALESCE(AVG(price), 0)
		FROM listings
		WHERE is_active = TRUE
		GROUP BY bt
		ORDER BY COUNT(*) DESC
	`)
	if err != nil {
		return stats, fmt.Errorf("get body types: %w", err)
	}
	defer btRows.Close()

	for btRows.Next() {
		var s models.BodyTypeStat
		if err := btRows.Scan(&s.Type, &s.Count, &s.AvgPrice); err != nil {
			return stats, fmt.Errorf("scan body type row: %w", err)
		}
		stats.BodyTypes = append(stats.BodyTypes, s)
	}
	if err := btRows.Err(); err != nil {
		return stats, fmt.Errorf("body type rows error: %w", err)
	}
	if stats.BodyTypes == nil {
		stats.BodyTypes = make([]models.BodyTypeStat, 0)
	}

	// Year distribution — only rows where year is known.
	yrRows, err := pool.Query(ctx, `
		SELECT year::int, COUNT(*)::int
		FROM listings
		WHERE is_active = TRUE AND year IS NOT NULL
		GROUP BY year
		ORDER BY year ASC
	`)
	if err != nil {
		return stats, fmt.Errorf("get year distribution: %w", err)
	}
	defer yrRows.Close()

	for yrRows.Next() {
		var y models.YearStat
		if err := yrRows.Scan(&y.Year, &y.Count); err != nil {
			return stats, fmt.Errorf("scan year row: %w", err)
		}
		stats.YearDistribution = append(stats.YearDistribution, y)
	}
	if err := yrRows.Err(); err != nil {
		return stats, fmt.Errorf("year rows error: %w", err)
	}
	if stats.YearDistribution == nil {
		stats.YearDistribution = make([]models.YearStat, 0)
	}

	return stats, nil
}

// strVal dereferences a *string, returning "" for nil pointers.
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

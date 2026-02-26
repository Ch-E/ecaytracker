package main

import (
	"context"
	"log"
	"os"

	"ecaytracker/backend/config"
	appdb "ecaytracker/backend/internal/db"
	"ecaytracker/backend/internal/scraper"
)

func main() {
	log.SetOutput(os.Stderr)

	cfg := config.Load()

	pool, err := appdb.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Database connected.")

	listings, err := scraper.Scrape()
	if err != nil {
		log.Fatalf("Scrape failed: %v", err)
	}
	if len(listings) == 0 {
		log.Fatal("No listings extracted — selectors may need updating or Cloudflare blocked the request.")
	}
	log.Printf("Scraped %d listing(s). Upserting to database…", len(listings))

	ctx := context.Background()
	var inserted, updated, priceChanged int

	for _, l := range listings {
		result, err := appdb.UpsertListing(ctx, pool, l)
		if err != nil {
			log.Printf("ERROR upserting %s (%s): %v", l.ExternalID, l.Title, err)
			continue
		}
		switch {
		case result.Inserted:
			inserted++
		case result.PriceChanged:
			priceChanged++
			updated++
		default:
			updated++
		}
	}

	log.Printf("Done — inserted: %d | updated: %d (price changed: %d) | errors: %d",
		inserted, updated, priceChanged, len(listings)-inserted-updated)
}

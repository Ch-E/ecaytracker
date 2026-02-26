# EcayTracker — Backend

Go backend for EcayTracker. Scrapes car listings from [ecaytrade.com](https://ecaytrade.com) and serves them via a REST API backed by Supabase PostgreSQL.

## Structure

```
backend/
  cmd/api/main.go          # Gin API server
  cmd/scraper/main.go      # Scraper → upserts to DB
  config/config.go         # Env var loader
  internal/
    api/
      handlers/            # health.go, listings.go
      router.go
    db/
      db.go                # pgxpool init
      queries.go           # upsert + fetch
    scraper/
      scraper.go           # rod + stealth browser
      parser.go            # DOM → Listing struct
  models/listing.go
  schema.sql
```

## Setup

### 1. Create the database

Run [schema.sql](./schema.sql) in the Supabase SQL Editor.

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env and set DATABASE_URL to your Supabase connection string
```

Get your connection string from **Supabase → Project Settings → Database → Connection string** (use the `postgresql://...` URI format).

### 3. Install Go dependencies

```bash
go mod tidy
```

## Running

### Scraper

Scrapes the first page of ecaytrade.com/autos-boats/autos and upserts listings into the database.

```bash
go run ./cmd/scraper

# Headless mode (no browser window)
HEADLESS=true go run ./cmd/scraper
```

### API Server

```bash
go run ./cmd/api
```

Endpoints:

| Method | Path             | Description                    |
|--------|------------------|--------------------------------|
| GET    | `/health`        | DB ping — 200 OK or 500        |
| GET    | `/api/listings`  | All active listings as JSON    |

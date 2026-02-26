-- Run this in Supabase SQL Editor (or any PostgreSQL client) before starting the app.

CREATE TABLE IF NOT EXISTS listings (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_id  TEXT UNIQUE NOT NULL,
  url          TEXT NOT NULL,
  title        TEXT NOT NULL,
  make         TEXT,
  model        TEXT,
  year         INTEGER,
  mileage      INTEGER,
  price        NUMERIC(10, 2),
  currency     TEXT DEFAULT 'KYD',
  condition    TEXT,
  transmission TEXT,
  fuel_type    TEXT,
  color        TEXT,
  description  TEXT,
  images       TEXT[],
  location     TEXT,
  seller_name  TEXT,
  is_active    BOOLEAN DEFAULT TRUE,
  first_seen   TIMESTAMPTZ DEFAULT NOW(),
  last_seen    TIMESTAMPTZ DEFAULT NOW(),
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_history (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  listing_id  UUID REFERENCES listings(id) ON DELETE CASCADE,
  price       NUMERIC(10, 2) NOT NULL,
  recorded_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index to speed up listing lookups
CREATE INDEX IF NOT EXISTS idx_listings_external_id ON listings(external_id);
CREATE INDEX IF NOT EXISTS idx_listings_is_active   ON listings(is_active);
CREATE INDEX IF NOT EXISTS idx_listings_make        ON listings(make);
CREATE INDEX IF NOT EXISTS idx_listings_created_at  ON listings(created_at DESC);

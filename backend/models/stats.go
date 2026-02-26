package models

// Stats holds pre-computed dashboard statistics.
type Stats struct {
	TotalListings int         `json:"total_listings"`
	AvgPrice      float64     `json:"avg_price"`
	MedianPrice   float64     `json:"median_price"`
	NewThisWeek   int         `json:"new_this_week"`
	TopBrands     []BrandStat `json:"top_brands"`
}

// BrandStat holds per-brand aggregates.
type BrandStat struct {
	Name     string  `json:"name"`
	Count    int     `json:"count"`
	AvgPrice float64 `json:"avg_price"`
}

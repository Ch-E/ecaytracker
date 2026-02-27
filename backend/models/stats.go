package models

// Stats holds pre-computed dashboard statistics.
type Stats struct {
	TotalListings    int            `json:"total_listings"`
	AvgPrice         float64        `json:"avg_price"`
	MedianPrice      float64        `json:"median_price"`
	NewThisWeek      int            `json:"new_this_week"`
	AvgMileage       float64        `json:"avg_mileage"`
	TopBrands        []BrandStat    `json:"top_brands"`
	BodyTypes        []BodyTypeStat `json:"body_types"`
	YearDistribution []YearStat     `json:"year_distribution"`
}

// BrandStat holds per-brand aggregates.
type BrandStat struct {
	Name     string  `json:"name"`
	Count    int     `json:"count"`
	AvgPrice float64 `json:"avg_price"`
}

// BodyTypeStat holds per-body-type aggregates.
type BodyTypeStat struct {
	Type     string  `json:"type"`
	Count    int     `json:"count"`
	AvgPrice float64 `json:"avg_price"`
}

// YearStat holds listing count for a model year.
type YearStat struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

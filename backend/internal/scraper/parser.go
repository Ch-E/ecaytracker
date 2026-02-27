// Package scraper provides functionality for parsing vehicle listing cards.
package scraper

import (
	"regexp"
	"strconv"
	"strings"

	"ecaytracker/backend/models"
)

var (
	// Matches prices like "CI$ 5,000", "KYD$ 6,000", "US$ 16,000", "CI$5000"
	priceRe = regexp.MustCompile(`(?i)(CI\$|KYD\$?|US\$?)\s*([\d,]+(?:\.\d+)?)`)

	// Matches a 4-digit year between 1900-2099
	yearRe = regexp.MustCompile(`\b((?:19|20)\d{2})\b`)

	// Matches the listing ID from the URL path
	idRe = regexp.MustCompile(`/advert/(\d+)$`)

	// Matches "On Island" / "Off Island" / "Grand Cayman" etc.
	locationRe = regexp.MustCompile(`(?i)(on island|off island|grand cayman|cayman brac|little cayman|george town|bodden town|west bay|north side|east end)`)

	// mileageRe matches a bare mileage figure inside a card snippet.
	mileageRe = regexp.MustCompile(`(?i)(over\s+[\d,]+|under\s+[\d,]+|[\d,]+\s*(?:km|miles|mi))`)

	// contaminatedRe detects a field value that contains another label's data,
	// e.g. "Mileage: 85600 km" appearing as the value for the "Drive" field.
	contaminatedRe = regexp.MustCompile(`(?i)^(?:mileage|odometer|body\s+type|drive|cylinders|fuel\s+type|transmission|steering|exterior\s+color|interior\s+color|doors|on\s+island|condition|year)\s*:`)

	// mileageDetailRe is used on full listing-detail pages.
	// Priority 1 – explicit label:  "Mileage: 45,000" / "Odometer: 45000 km"
	// Priority 2 – value + unit:    "45,000 km" / "100000 miles"
	// Priority 3 – approximate:     "Over 100,000" / "Under 50,000"
	mileageLabelRe  = regexp.MustCompile(`(?i)(?:mileage|odometer|miles|km|kilometers)\s*[:\-–]?\s*([\d,]+)\s*(?:km|miles|mi)?`)
	mileageValueRe  = regexp.MustCompile(`(?i)([\d,]+)\s*(?:km|kilometers|miles|mi)\b`)
	mileageApproxRe = regexp.MustCompile(`(?i)(?:over|under|approx\.?)\s+([\d,]+)`)

	// digitsRe extracts the first run of digits (and commas) from a string.
	digitsRe = regexp.MustCompile(`[\d,]+`)
)

// knownMakes is a rough list of common makes to help split make vs model.
var knownMakes = []string{
	"Acura", "Alfa Romeo", "Aston Martin", "Audi", "Bentley", "BMW", "Bugatti",
	"Buick", "Cadillac", "Chevrolet", "Chrysler", "Citroën", "Dodge", "Ferrari",
	"Fiat", "Ford", "Genesis", "GMC", "Honda", "Hyundai", "Infiniti", "Jaguar",
	"Jeep", "Kia", "Lamborghini", "Land Rover", "Lexus", "Lincoln", "Lotus",
	"Maserati", "Mazda", "McLaren", "Mercedes", "Mercedes-Benz", "MINI", "Mitsubishi",
	"Nissan", "Peugeot", "Pontiac", "Porsche", "Ram", "Rolls-Royce", "Subaru",
	"Suzuki", "Tesla", "Toyota", "Volkswagen", "Volvo",
}

// ParseCard parses a raw listing card into a Listing struct.
// cardText is the full innerText of the anchor element.
func ParseCard(cardText, rawURL, imgURL string) models.Listing {
	l := models.Listing{
		URL:      rawURL,
		IsActive: true,
	}
	if imgURL != "" {
		l.Images = []string{imgURL}
	}

	// Extract external ID from URL
	if m := idRe.FindStringSubmatch(rawURL); len(m) == 2 {
		l.ExternalID = m[1]
	}

	// Extract price + currency
	if m := priceRe.FindStringSubmatch(cardText); len(m) == 3 {
		currency := normaliseCurrency(m[1])
		priceStr := strings.ReplaceAll(m[2], ",", "")
		if v, err := strconv.ParseFloat(priceStr, 64); err == nil {
			l.Price = v
			l.Currency = currency
		}
	}

	// Extract year
	years := yearRe.FindAllString(cardText, -1)
	if len(years) > 0 {
		// Prefer the last found year (sometimes the title starts with year)
		if v, err := strconv.Atoi(years[len(years)-1]); err == nil {
			l.Year = &v
		}
	}

	// Extract location
	if m := locationRe.FindString(cardText); m != "" {
		l.Location = normaliseTitle(m)
	}

	// Extract mileage from card text (also caught by ParseMileage on detail page).
	l.Mileage = ParseMileage(cardText)

	// Build title: take the first non-empty line that isn't just a price
	l.Title = extractTitle(cardText)

	// Split make + model from title
	l.Make, l.Model = splitMakeModel(l.Title)

	return l
}

// extractTitle picks the most meaningful line from the card text as the listing title.
func extractTitle(text string) string {
	// Normalise whitespace
	text = strings.ReplaceAll(text, "\t", " ")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Skip lines that are purely a price
		if priceRe.MatchString(line) && len(line) < 20 {
			continue
		}
		// Skip lines that are purely details ("Automatic · 2018 · On Island")
		if strings.Contains(line, "·") {
			continue
		}
		return line
	}
	// Fallback: return first 60 chars of cleaned text
	clean := strings.Join(strings.Fields(text), " ")
	if len(clean) > 60 {
		return clean[:60]
	}
	return clean
}

// splitMakeModel attempts to extract make and model from a title string.
// e.g. "2018 Toyota Camry SE" -> ("Toyota", "Camry SE")
func splitMakeModel(title string) (make_, model string) {
	// Strip leading year if present
	stripped := yearRe.ReplaceAllString(title, "")
	stripped = strings.TrimSpace(stripped)

	titleUpper := strings.ToUpper(stripped)
	for _, mk := range knownMakes {
		if strings.HasPrefix(titleUpper, strings.ToUpper(mk)) {
			make_ = mk
			model = strings.TrimSpace(stripped[len(mk):])
			return
		}
	}

	// Fallback: first word is make, rest is model
	parts := strings.Fields(stripped)
	if len(parts) == 0 {
		return "", title
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func normaliseCurrency(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch {
	case strings.HasPrefix(s, "CI"):
		return "KYD" // CI$ == Cayman Islands Dollar == KYD
	case strings.HasPrefix(s, "KYD"):
		return "KYD"
	case strings.HasPrefix(s, "US"):
		return "USD"
	default:
		return s
	}
}

func normaliseTitle(s string) string {
	return strings.Title(strings.ToLower(s)) //nolint:staticcheck
}

// ParseMileage extracts a mileage integer from arbitrary text (card snippet or
// full detail-page body). Returns nil when no mileage can be determined.
// Three patterns are tried in priority order:
//  1. Labeled field  – "Mileage: 45,000" / "Odometer: 100,000 km"
//  2. Value + unit   – "45,000 km" / "100,000 miles"
//  3. Approximate    – "Over 100,000" / "Under 50,000"
func ParseMileage(text string) *int {
	try := func(re *regexp.Regexp) *int {
		m := re.FindStringSubmatch(text)
		if len(m) < 2 {
			return nil
		}
		digits := strings.ReplaceAll(digitsRe.FindString(m[1]), ",", "")
		v, err := strconv.Atoi(digits)
		if err != nil || v < 100 || v > 2_000_000 {
			return nil
		}
		return &v
	}

	if v := try(mileageLabelRe); v != nil {
		return v
	}
	if v := try(mileageValueRe); v != nil {
		return v
	}
	return try(mileageApproxRe)
}

// ApplyDetailFields merges a structured key→value map (extracted from an "Ad
// Details" section) and the full page text into a Listing. Existing non-zero
// values are not overwritten so card-level data always wins.
func ApplyDetailFields(fields map[string]string, fullText string, l *models.Listing) {
	get := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := fields[k]; ok && v != "" {
				v = strings.TrimSpace(v)
				// Reject values that look like another field's label+colon pair,
				// indicating the DOM walker grabbed the wrong neighbouring element.
				if contaminatedRe.MatchString(v) {
					continue
				}
				return v
			}
		}
		return ""
	}

	if l.Mileage == nil {
		// Try the structured fields map first (most reliable).
		if raw := get("mileage"); raw != "" {
			l.Mileage = parseMileageValue(raw)
		}
		// Fall back to full-text regex when the map had nothing.
		if l.Mileage == nil && fullText != "" {
			l.Mileage = ParseMileage(fullText)
		}
	}

	if l.Condition == "" {
		l.Condition = get("condition")
	}
	if l.Transmission == "" {
		l.Transmission = get("transmission")
	}
	if l.FuelType == "" {
		l.FuelType = get("fuel type", "fuel_type")
	}
	if l.Color == "" {
		l.Color = get("exterior color", "color")
	}
	if l.BodyType == "" {
		l.BodyType = get("body type", "body_type")
	}
	if l.Drive == "" {
		l.Drive = get("drive")
	}
	if l.Cylinders == "" {
		l.Cylinders = get("cylinders")
	}
	if l.Steering == "" {
		l.Steering = get("steering")
	}
	if l.InteriorColor == "" {
		l.InteriorColor = get("interior color", "interior_color")
	}
	if l.Doors == "" {
		l.Doors = get("doors")
	}
	if l.OnIsland == nil {
		if raw := get("on island", "on_island"); raw != "" {
			yes := strings.EqualFold(raw, "yes") || strings.EqualFold(raw, "true")
			l.OnIsland = &yes
		}
	}
	// Year from detail page (only override if not already set from card text).
	if l.Year == nil {
		if raw := get("year"); raw != "" {
			if v, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && v >= 1900 && v <= 2100 {
				l.Year = &v
			}
		}
	}
}

// parseMileageValue converts a raw mileage string (e.g. "Over 100,000" or
// "48,050") into an integer, applying the same sanity bounds as ParseMileage.
func parseMileageValue(raw string) *int {
	// Strip leading "over"/"under"/"approx" qualifier.
	raw = regexp.MustCompile(`(?i)^(over|under|approx\.?)\s+`).ReplaceAllString(raw, "")
	digits := strings.ReplaceAll(digitsRe.FindString(raw), ",", "")
	v, err := strconv.Atoi(digits)
	if err != nil || v < 100 || v > 2_000_000 {
		return nil
	}
	return &v
}

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

	// Matches mileage description
	mileageRe = regexp.MustCompile(`(?i)(over\s+[\d,]+|under\s+[\d,]+|[\d,]+\s*(?:km|miles|mi))`)
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

	// Extract mileage if present in card text (strip non-numeric prefix like "over"/"under").
	if m := mileageRe.FindString(cardText); m != "" {
		// Extract only the first digit sequence from matches like "Over 100,000"
		digits := regexp.MustCompile(`[\d,]+`).FindString(m)
		digits = strings.ReplaceAll(digits, ",", "")
		if v, err := strconv.Atoi(digits); err == nil {
			l.Mileage = &v
		}
	}

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

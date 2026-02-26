package scraper

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"

	"ecaytracker/backend/models"
)

const (
	listingsURL = "https://ecaytrade.com/autos-boats/autos"

	// cardSelector matches all advert listing anchor elements on the page.
	cardSelector = `a[href*="/advert/"]`
)

// Scrape launches a browser, navigates to the ecaytrade autos listing page,
// and returns all parsed listings found on the first page.
func Scrape() ([]models.Listing, error) {
	headless := strings.EqualFold(os.Getenv("HEADLESS"), "true")
	log.Printf("Launching browser (headless=%v)...", headless)

	// Leakless(false) disables the leakless.exe helper binary that Windows Defender
	// incorrectly flags as malware (known false positive with go-rod on Windows).
	l := launcher.New().Headless(headless).Leakless(false)
	u, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("launcher: %w", err)
	}

	browser := rod.New().ControlURL(u).MustConnect()
	defer func() {
		if cerr := browser.Close(); cerr != nil {
			log.Printf("browser.Close: %v", cerr)
		}
	}()

	// stealth.Page creates the page AND injects all anti-detection scripts
	// before any navigation -- this is the correct way to use go-rod/stealth.
	page, err := stealth.Page(browser)
	if err != nil {
		return nil, fmt.Errorf("stealth.Page: %w", err)
	}

	log.Printf("Navigating to %s ...", listingsURL)
	if err := page.Navigate(listingsURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}

	// Wait for the document to finish loading.
	page.MustWaitLoad()

	// Additionally wait for at least one listing card anchor to appear.
	log.Println("Waiting for listing cards in DOM...")
	if err := waitForSelector(page, cardSelector, 30*time.Second); err != nil {
		html, _ := page.HTML()
		log.Printf("Page HTML dump (first 3000 chars):\n%s", truncate(html, 3000))
		return nil, fmt.Errorf("listing cards never appeared: %w", err)
	}

	// Random human-like pause (2-3 s) before extracting.
	delay := time.Duration(2000+rand.Intn(1000)) * time.Millisecond
	log.Printf("Sleeping %v (human delay)...", delay)
	time.Sleep(delay)

	// Extract raw card data via a single JS round-trip.
	log.Println("Extracting cards via JS eval...")
	cards, err := extractCards(page)
	if err != nil {
		return nil, fmt.Errorf("extractCards: %w", err)
	}
	log.Printf("Found %d raw card(s)", len(cards))

	var listings []models.Listing
	for i, c := range cards {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[card %d] panic: %v", i, r)
				}
			}()
			listing := ParseCard(c.Text, c.URL, c.ImgSrc)
			if listing.Title == "" && listing.Price == 0 {
				log.Printf("[card %d] skipped -- empty title+price (url=%s)", i, c.URL)
				return
			}
			if listing.Year == nil {
				log.Printf("[card %d] skipped — no year (likely a part/accessory): %s", i, listing.Title)
				return
			}
			listings = append(listings, listing)
		}()
	}

	return listings, nil
}

// rawCard holds the DOM data extracted by the JS snippet.
type rawCard struct {
	URL    string
	Text   string
	ImgSrc string
}

// extractCards runs JavaScript on the page to collect all advert links.
func extractCards(page *rod.Page) ([]rawCard, error) {
	result, err := page.Eval(`() => {
const seen = new Set();
return [...document.querySelectorAll('a[href*="/advert/"]')]
.filter(a => {
if (!/\/advert\/\d+$/.test(a.href)) return false;
if (seen.has(a.href)) return false;
seen.add(a.href);
return true;
})
.map(a => {
const img = a.querySelector('img');
return {
url:    a.href,
text:   a.innerText.trim(),
imgSrc: img ? img.src : ''
};
});
}`)
	if err != nil {
		return nil, fmt.Errorf("JS eval: %w", err)
	}

	raw := result.Value
	if raw.Nil() {
		return nil, fmt.Errorf("JS eval returned null/undefined")
	}

	var cards []rawCard
	for _, item := range raw.Arr() {
		obj := item.Map()
		cards = append(cards, rawCard{
			URL:    obj["url"].Str(),
			Text:   obj["text"].Str(),
			ImgSrc: obj["imgSrc"].Str(),
		})
	}
	return cards, nil
}

// waitForSelector polls until the selector matches at least one element.
func waitForSelector(page *rod.Page, selector string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		el, err := page.Element(selector)
		if err == nil && el != nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("selector %q not found within %v", selector, timeout)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// ---------------------------------------------------------------------------
// PROXY FALLBACK NOTE
// ---------------------------------------------------------------------------
// If Cloudflare blocks the request even with stealth enabled, route traffic
// through a residential proxy service such as ScrapingBee:
//
//   l = launcher.New().
//       Headless(true).
//       Proxy("http://YOUR_SCRAPINGBEE_PROXY_HOST:PORT").
//       Set("--ignore-certificate-errors")
//
// ScrapingBee also offers a simple HTTP API (no browser required):
//   https://www.scrapingbee.com/documentation/
// ---------------------------------------------------------------------------

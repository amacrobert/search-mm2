package scraper

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"search-mm2/backend/internal/models"
)

type browserProfile struct {
	UserAgent    string
	SecChUa      string
	Platform     string
}

var browserProfiles = []browserProfile{
	{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		SecChUa:   `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`,
		Platform:  `"macOS"`,
	},
	{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		SecChUa:   `"Chromium";v="123", "Not/A)Brand";v="8", "Google Chrome";v="123"`,
		Platform:  `"Windows"`,
	},
	{
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		SecChUa:   `"Chromium";v="121", "Not A(Brand";v="99", "Google Chrome";v="121"`,
		Platform:  `"Linux"`,
	},
}

var errEmptyURL = fmt.Errorf("search URL is empty")

func ScrapeLoopNet(search *models.Search) ([]models.Property, error) {
	if strings.TrimSpace(search.URL) == "" {
		return nil, errEmptyURL
	}

	var properties []models.Property

	profile := browserProfiles[rand.IntN(len(browserProfiles))]

	c := colly.NewCollector(
		colly.UserAgent(profile.UserAgent),
		colly.AllowURLRevisit(),
	)

	jar, _ := cookiejar.New(nil)
	c.SetCookieJar(jar)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*loopnet.*",
		Delay:       3 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Sec-Ch-Ua", profile.SecChUa)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", profile.Platform)
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	c.OnHTML("article.placard", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("h4.placard-header-title, .placard-header-title"))
		address := strings.TrimSpace(e.ChildText(".placard-header-subtitle"))
		listingURL := e.ChildAttr("a.placard-header-link, a.placard-pseudo-link", "href")
		if listingURL != "" && !strings.HasPrefix(listingURL, "http") {
			listingURL = "https://www.loopnet.com" + listingURL
		}
		imageURL := e.ChildAttr("img.placard-photo, .placard-photo img", "src")

		priceText := strings.TrimSpace(e.ChildText(".placard-header-price, .data-points-1 li:first-child"))
		sizeText := strings.TrimSpace(e.ChildText(".data-points-1 li:nth-child(2), .placard-header-size"))
		propType := strings.TrimSpace(e.ChildText(".data-points-2 li:first-child, .placard-header-type"))

		externalID := e.Attr("data-id")
		if externalID == "" {
			externalID = e.Attr("id")
		}
		if externalID == "" && listingURL != "" {
			parts := strings.Split(strings.TrimRight(listingURL, "/"), "/")
			if len(parts) > 0 {
				externalID = parts[len(parts)-1]
			}
		}

		if name == "" && address == "" {
			return
		}

		city, state, zip := parseAddress(address)

		p := models.Property{
			ExternalID:   externalID,
			Name:         name,
			Address:      address,
			City:         city,
			State:        state,
			Zip:          zip,
			PropertyType: propType,
			Price:        parsePrice(priceText),
			SizeSqFt:     parseSize(sizeText),
			URL:          listingURL,
			ImageURL:     imageURL,
			ScrapedAt:    time.Now(),
		}
		properties = append(properties, p)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("scraper: request error (status %d): %s %v", r.StatusCode, r.Request.URL, err)
	})

	const maxAttempts = 3
	backoffs := []time.Duration{5 * time.Second, 20 * time.Second}

	for attempt := range maxAttempts {
		properties = nil
		log.Printf("scraper: fetching %s (attempt %d/%d)", search.URL, attempt+1, maxAttempts)
		if err := c.Visit(search.URL); err != nil {
			if attempt == maxAttempts-1 {
				return nil, fmt.Errorf("visit %s: %w", search.URL, err)
			}
			log.Printf("scraper: attempt %d failed: %v, retrying...", attempt+1, err)
			time.Sleep(backoffs[attempt])
			continue
		}
		c.Wait()
		if len(properties) > 0 {
			break
		}
		if attempt < maxAttempts-1 {
			log.Printf("scraper: attempt %d returned 0 properties, retrying...", attempt+1)
			time.Sleep(backoffs[attempt])
		}
	}

	return properties, nil
}


func parseAddress(addr string) (city, state, zip string) {
	parts := strings.Split(addr, ",")
	if len(parts) >= 2 {
		city = strings.TrimSpace(parts[len(parts)-2])
		stateZip := strings.TrimSpace(parts[len(parts)-1])
		szParts := strings.Fields(stateZip)
		if len(szParts) >= 1 {
			state = szParts[0]
		}
		if len(szParts) >= 2 {
			zip = szParts[1]
		}
	}
	return
}

func parsePrice(s string) *float64 {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var val float64
	if _, err := fmt.Sscanf(s, "%f", &val); err == nil {
		return &val
	}
	return nil
}

func parseSize(s string) *int {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "SF", "")
	s = strings.ReplaceAll(s, "sf", "")
	s = strings.ReplaceAll(s, "sqft", "")
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err == nil {
		return &val
	}
	return nil
}

package leipzigde

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/havocked/leipzig-cli/internal/model"
)

const baseURL = "https://www.leipzig.de"

var topicToCategory = map[string]string{
	"Konzert":            model.CategoryConcert,
	"Klassik":            model.CategoryConcert,
	"Jazz und Blues":     model.CategoryConcert,
	"Kabarett":           model.CategoryTheater,
	"Oper und Operette":  model.CategoryTheater,
	"Führungen":          model.CategoryCulture,
	"Kurse und Treffs":   model.CategoryCulture,
	"Mitmach-Angebot":    model.CategoryCulture,
	"Beratung":           model.CategoryOther,
	"Ausstellungen":      model.CategoryExhibition,
	"Lesung":             model.CategoryCulture,
	"Kinder & Jugendliche": model.CategoryFamily,
	"Freizeit":           model.CategoryOther,
	"Bühne":              model.CategoryTheater,
	"Sport":              model.CategorySport,
	"Märkte":             model.CategoryMarket,
}

type Source struct {
	client *http.Client
}

func New() *Source {
	return &Source{client: &http.Client{Timeout: 30 * time.Second}}
}

func (s *Source) ID() string { return "leipzig.de" }

func (s *Source) Fetch(ctx context.Context, from, to time.Time) ([]model.Event, error) {
	urls := []string{
		baseURL + "/kultur-und-freizeit/veranstaltungen/termine-heute",
	}

	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Friday || weekday == time.Saturday || weekday == time.Sunday {
		urls = append(urls, baseURL+"/kultur-und-freizeit/veranstaltungen/termine-dieses-wochenende")
	}

	var allEvents []model.Event
	seen := make(map[string]bool)

	for i, u := range urls {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		events, err := s.fetchPage(ctx, u)
		if err != nil {
			fmt.Fprintf(nil, "")
			return nil, fmt.Errorf("fetch %s: %w", u, err)
		}

		for _, e := range events {
			key := e.Name + "|" + e.StartTime.String() + "|" + e.Venue
			if !seen[key] {
				seen[key] = true
				if !from.IsZero() && e.StartTime.Before(from) {
					continue
				}
				if !to.IsZero() && e.StartTime.After(to) {
					continue
				}
				allEvents = append(allEvents, e)
			}
		}
	}

	return allEvents, nil
}

func (s *Source) fetchPage(ctx context.Context, url string) ([]model.Event, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "leipzig-cli/1.0")
	req.Header.Set("Accept-Language", "de-DE,de;q=0.9")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []model.Event

	doc.Find("li[data-event] article.event-card, li[data-event] article.card").Each(func(_ int, card *goquery.Selection) {
		e := model.Event{Source: "leipzig.de"}

		// Name from h3
		e.Name = strings.TrimSpace(card.Find("h3").First().Text())

		// URL from link
		if href, ok := card.Find("a[href]").First().Attr("href"); ok {
			if strings.HasPrefix(href, "/") {
				e.URL = baseURL + href
			} else if strings.HasPrefix(href, "http") {
				e.URL = href
			}
		}

		// Image URL
		if src, ok := card.Find("img").First().Attr("src"); ok {
			if strings.HasPrefix(src, "/") {
				e.ImageURL = baseURL + src
			} else {
				e.ImageURL = src
			}
		}

		// Parse icon-based fields: each is a span.icon-text containing span.icon + span with value
		card.Find("span.icon-text").Each(func(_ int, iconText *goquery.Selection) {
			icon := strings.TrimSpace(iconText.Find("span.icon").First().Text())
			// Value is in the sibling span (not the icon span)
			var value string
			iconText.Find("span").Each(func(_ int, s *goquery.Selection) {
				if !s.HasClass("icon") {
					t := strings.TrimSpace(s.Text())
					if t != "" {
						value = t
					}
				}
			})

			switch icon {
			case "event":
				e.StartTime, e.EndTime = parseDateTime(value)
			case "location_on":
				e.Venue = value
			case "topic":
				e.Category = mapTopicToCategory(value)
			}
		})

		if e.Category == "" || e.Category == model.CategoryOther {
			e.Category = model.InferCategory(e.Name, e.Venue)
		}

		if e.Name != "" {
			events = append(events, e)
		}
	})

	return events, nil
}

func parseDateTime(raw string) (start, end time.Time) {
	loc, _ := time.LoadLocation("Europe/Berlin")
	raw = strings.ReplaceAll(raw, "\u00b7", "·")
	raw = strings.ReplaceAll(raw, "\u2013", "–")
	raw = strings.TrimSpace(raw)

	// Format a: "DD.MM.YYYY - DD.MM.YYYY" (date range)
	if strings.Contains(raw, " - ") && !strings.Contains(raw, "·") {
		parts := strings.SplitN(raw, " - ", 2)
		start, _ = time.ParseInLocation("02.01.2006", strings.TrimSpace(parts[0]), loc)
		end, _ = time.ParseInLocation("02.01.2006", strings.TrimSpace(parts[1]), loc)
		return
	}

	// Format b/c: contains "·"
	if strings.Contains(raw, "·") {
		parts := strings.SplitN(raw, "·", 2)
		dateStr := strings.TrimSpace(parts[0])
		timeStr := strings.TrimSpace(parts[1])
		timeStr = strings.ReplaceAll(timeStr, "Uhr", "")
		timeStr = strings.TrimSpace(timeStr)

		// Format c: "HH:MM – HH:MM"
		if strings.Contains(timeStr, "–") {
			tParts := strings.SplitN(timeStr, "–", 2)
			startTime := strings.TrimSpace(tParts[0])
			endTime := strings.TrimSpace(tParts[1])
			start, _ = time.ParseInLocation("02.01.2006 15:04", dateStr+" "+startTime, loc)
			end, _ = time.ParseInLocation("02.01.2006 15:04", dateStr+" "+endTime, loc)
			return
		}

		// Format b: "HH:MM"
		start, _ = time.ParseInLocation("02.01.2006 15:04", dateStr+" "+timeStr, loc)
		return
	}

	// Fallback: just a date
	start, _ = time.ParseInLocation("02.01.2006", raw, loc)
	return
}

func mapTopicToCategory(topic string) string {
	topic = strings.TrimSpace(topic)

	// Try exact match first
	if cat, ok := topicToCategory[topic]; ok {
		return cat
	}

	// Try compound topics (split by " · ")
	parts := strings.Split(topic, " · ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if cat, ok := topicToCategory[p]; ok {
			return cat
		}
	}

	return model.CategoryOther
}

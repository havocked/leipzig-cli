package prinzde

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/havocked/leipzig-cli/internal/model"
)

const baseURL = "https://prinz.de/leipzig/events/"

var categoryTextMap = map[string]string{
	"KONZERTE & LIVEMUSIK": model.CategoryConcert,
	"BÜHNE":                model.CategoryTheater,
	"AUSSTELLUNGEN":        model.CategoryExhibition,
	"KINDER & FAMILIE":     model.CategoryFamily,
	"SPORT":                model.CategorySport,
	"FÜHRUNGEN":            model.CategoryCulture,
	"STADTLEBEN":           model.CategoryCulture,
	"SPECIAL EVENTS":       model.CategoryCulture,
	"ESSEN & TRINKEN":      model.CategoryOther,
}

type Source struct{}

func New() *Source { return &Source{} }

func (s *Source) ID() string { return "prinz.de" }

func (s *Source) Fetch(ctx context.Context, from, to time.Time) ([]model.Event, error) {
	url := pickURL(from, to)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("prinzde: create request: %w", err)
	}
	req.Header.Set("User-Agent", "leipzig-cli/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prinzde: fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prinzde: %s returned %d", url, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("prinzde: parse HTML: %w", err)
	}

	loc, _ := time.LoadLocation("Europe/Berlin")
	var events []model.Event

	doc.Find("article.event-teaser").Each(func(_ int, card *goquery.Selection) {
		titleEl := card.Find("h3.event-teaser-title a")
		name := strings.TrimSpace(titleEl.Text())
		if name == "" {
			return
		}

		href, _ := titleEl.Attr("href")
		eventURL := href
		if eventURL != "" && !strings.HasPrefix(eventURL, "http") {
			eventURL = "https://prinz.de" + eventURL
		}

		// Category
		catText := strings.TrimSpace(card.Find(".event-teaser-category").Text())
		category := mapCategory(catText)

		// Date
		dateText := strings.TrimSpace(card.Find(".text-primary.text-sm-end").Text())
		eventDate := parseDate(dateText, loc)

		// Time + Venue from .event-teaser-meta
		meta := card.Find(".event-teaser-meta")
		timeText := strings.TrimSpace(meta.Find("span.fw-bold").Text())
		venue := strings.TrimSpace(meta.Find("span.text-uppercase").Text())

		startTime := applyTime(eventDate, timeText, loc)

		// Image
		imageURL, _ := card.Find(".teaser-thumbnail img").Attr("src")

		events = append(events, model.Event{
			Name:      name,
			StartTime: startTime,
			Venue:     venue,
			Category:  category,
			URL:       eventURL,
			ImageURL:  imageURL,
			Source:    "prinz.de",
		})
	})

	return events, nil
}

func pickURL(from, to time.Time) string {
	loc, _ := time.LoadLocation("Europe/Berlin")
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	tomorrow := today.Add(24 * time.Hour)
	duration := to.Sub(from)

	if from.After(today) && from.Before(tomorrow.Add(24*time.Hour)) && duration <= 24*time.Hour {
		if from.After(today.Add(23 * time.Hour)) {
			return baseURL + "morgen/"
		}
	}

	if duration > 3*24*time.Hour {
		return baseURL + "7-tage/"
	}
	if duration > 24*time.Hour {
		return baseURL + "wochenende/"
	}
	if from.After(today.Add(23 * time.Hour)) {
		return baseURL + "morgen/"
	}
	return baseURL
}

func mapCategory(text string) string {
	parts := strings.Split(text, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if cat, ok := categoryTextMap[p]; ok && cat != model.CategoryOther {
			return cat
		}
	}
	// fallback: check first part
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if cat, ok := categoryTextMap[p]; ok {
			return cat
		}
	}
	return model.CategoryOther
}

func parseDate(text string, loc *time.Location) time.Time {
	// Format: "Mo. 16.02.26"
	text = strings.TrimSpace(text)
	parts := strings.Fields(text)
	if len(parts) < 2 {
		return time.Time{}
	}
	datePart := parts[len(parts)-1] // last token is DD.MM.YY
	t, err := time.ParseInLocation("02.01.06", datePart, loc)
	if err != nil {
		return time.Time{}
	}
	return t
}

func applyTime(date time.Time, timeText string, loc *time.Location) time.Time {
	if date.IsZero() {
		return date
	}
	timeText = strings.TrimSpace(timeText)
	if timeText == "" || strings.EqualFold(timeText, "ganztägig") {
		return date
	}
	t, err := time.ParseInLocation("15:04", timeText, loc)
	if err != nil {
		return date
	}
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, loc)
}

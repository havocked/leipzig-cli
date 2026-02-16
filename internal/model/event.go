package model

import (
	"fmt"
	"net/url"
	"time"
)

const (
	CategoryConcert    = "concert"
	CategoryTheater    = "theater"
	CategoryExhibition = "exhibition"
	CategoryFamily     = "family"
	CategoryMarket     = "market"
	CategorySport      = "sport"
	CategoryCulture    = "culture"
	CategoryNightlife  = "nightlife"
	CategoryOther      = "other"
)

var CategoryEmoji = map[string]string{
	CategoryConcert:    "🎵",
	CategoryTheater:    "🎭",
	CategoryExhibition: "🖼️",
	CategoryFamily:     "👨‍👩‍👧‍👦",
	CategoryMarket:     "🛍️",
	CategorySport:      "⚽",
	CategoryCulture:    "📚",
	CategoryNightlife:  "🌙",
	CategoryOther:      "📌",
}

type Event struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime,omitzero"`
	Venue       string    `json:"venue,omitempty"`
	Address     string    `json:"address,omitempty"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags,omitempty"`
	Price       string    `json:"price,omitempty"`
	URL         string    `json:"url,omitempty"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	MapURL      string    `json:"mapUrl,omitempty"`
	Source      string    `json:"source"`
}

// MapsURL returns a Google Maps search URL for the venue.
func (e Event) MapsURL() string {
	if e.Venue == "" {
		return ""
	}
	q := e.Venue + ", Leipzig"
	if e.Address != "" {
		q = e.Address + ", Leipzig"
	}
	return "https://maps.google.com/?q=" + url.QueryEscape(q)
}

func (e Event) String() string {
	emoji := CategoryEmoji[e.Category]
	if emoji == "" {
		emoji = "📌"
	}
	t := e.StartTime.Format("Mon 02 Jan 15:04")
	venue := ""
	if e.Venue != "" {
		venue = " @ " + e.Venue
	}
	price := ""
	if e.Price != "" {
		price = fmt.Sprintf(" (%s)", e.Price)
	}
	return fmt.Sprintf("%s | %s %s%s%s", t, emoji, e.Name, venue, price)
}

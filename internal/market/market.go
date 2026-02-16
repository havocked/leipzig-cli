package market

import (
	"fmt"
	"net/url"
	"time"
)

type Schedule struct {
	Day   time.Weekday
	Open  string
	Close string
}

type Market struct {
	Name     string     `json:"name"`
	Private  bool       `json:"private,omitempty"`
	Schedules []Schedule `json:"-"`
	Notes    string     `json:"notes,omitempty"`
	MapURL   string     `json:"map_url"`
}

type MarketDay struct {
	Name    string `json:"name"`
	Open    string `json:"open"`
	Close   string `json:"close"`
	Notes   string `json:"notes,omitempty"`
	Private bool   `json:"private,omitempty"`
	MapURL  string `json:"map_url"`
}

func mapURL(location string) string {
	return "https://maps.google.com/?q=" + url.QueryEscape(location+", Leipzig")
}

func m(name string, private bool, schedules []Schedule, notes string) Market {
	return Market{
		Name:      name,
		Private:   private,
		Schedules: schedules,
		Notes:     notes,
		MapURL:    mapURL(name),
	}
}

var Markets = []Market{
	m("Innenstadt (Marktplatz)", false, []Schedule{
		{time.Tuesday, "09:00", "17:00"},
		{time.Friday, "09:00", "17:00"},
	}, ""),
	m("Bayrischer Platz", false, []Schedule{
		{time.Wednesday, "09:00", "17:00"},
		{time.Friday, "09:00", "17:00"},
	}, ""),
	m("Lindenauer Markt", false, []Schedule{
		{time.Wednesday, "09:00", "16:00"},
		{time.Friday, "09:00", "16:00"},
	}, ""),
	m("Gohlis-Park", false, []Schedule{
		{time.Tuesday, "09:00", "16:00"},
		{time.Thursday, "09:00", "15:00"},
	}, "Thu until 15:00"),
	m("Gohlis-Arkaden", false, []Schedule{
		{time.Wednesday, "09:00", "15:00"},
	}, ""),
	m("Lößnig", false, []Schedule{
		{time.Thursday, "09:00", "14:00"},
		{time.Saturday, "08:30", "12:00"},
	}, "Sat 08:30-12:00"),
	m("Grünau WK 4", false, []Schedule{
		{time.Tuesday, "09:00", "14:00"},
		{time.Thursday, "09:00", "14:00"},
	}, ""),
	m("Grünau WK 2", false, []Schedule{
		{time.Friday, "09:00", "12:00"},
	}, ""),
	m("Grünau WK 7", false, []Schedule{
		{time.Wednesday, "09:00", "12:00"},
	}, ""),
	m("Paunsdorf", false, []Schedule{
		{time.Thursday, "08:30", "13:30"},
	}, ""),
	m("Torgauer Platz", false, []Schedule{
		{time.Thursday, "09:00", "14:00"},
	}, ""),
	m("Richard-Wagner-Platz", false, []Schedule{
		{time.Saturday, "10:00", "16:00"},
	}, ""),
	m("Liebertwolkwitz", false, []Schedule{
		{time.Friday, "08:00", "13:00"},
	}, ""),
	m("Wiederitzsch", false, []Schedule{
		{time.Thursday, "09:00", "12:00"},
	}, ""),
	m("Sportforum", true, []Schedule{
		{time.Saturday, "09:00", "16:00"},
	}, ""),
	m("Plagwitzer Markthalle", true, []Schedule{
		{time.Saturday, "09:00", "14:00"},
	}, ""),
}

// ForDay returns markets open on the given weekday.
func ForDay(day time.Weekday) []MarketDay {
	var result []MarketDay
	for _, mk := range Markets {
		for _, s := range mk.Schedules {
			if s.Day == day {
				name := mk.Name
				if mk.Private {
					name += " (private)"
				}
				result = append(result, MarketDay{
					Name:    name,
					Open:    s.Open,
					Close:   s.Close,
					Notes:   mk.Notes,
					Private: mk.Private,
					MapURL:  mk.MapURL,
				})
			}
		}
	}
	return result
}

// AllByDay returns all markets grouped by weekday (Mon-Sun).
func AllByDay() map[time.Weekday][]MarketDay {
	days := map[time.Weekday][]MarketDay{}
	for d := time.Monday; d <= time.Saturday; d++ {
		if list := ForDay(d); len(list) > 0 {
			days[d] = list
		}
	}
	// Sunday
	if list := ForDay(time.Sunday); len(list) > 0 {
		days[time.Sunday] = list
	}
	return days
}

// FormatTime returns "HH:MM-HH:MM"
func FormatTime(open, close string) string {
	return fmt.Sprintf("%s-%s", open, close)
}

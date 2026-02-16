package engine

import (
	"strings"
	"time"

	"github.com/havocked/leipzig-cli/internal/model"
)

type FilterOptions struct {
	Category string
	Search   string
	Tags     []string
	Free     bool
	From     time.Time
	To       time.Time
	Limit    int
}

func Filter(events []model.Event, opts FilterOptions) []model.Event {
	var result []model.Event

	for _, e := range events {
		if opts.Category != "" {
			cats := strings.Split(opts.Category, ",")
			match := false
			for _, c := range cats {
				if strings.EqualFold(e.Category, strings.TrimSpace(c)) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		if !opts.From.IsZero() && e.StartTime.Before(opts.From) {
			continue
		}
		if !opts.To.IsZero() && e.StartTime.After(opts.To) {
			continue
		}

		if opts.Search != "" {
			q := strings.ToLower(opts.Search)
			if !strings.Contains(strings.ToLower(e.Name), q) &&
				!strings.Contains(strings.ToLower(e.Venue), q) {
				continue
			}
		}

		if opts.Free {
			p := strings.ToLower(e.Price)
			if p != "free" && p != "frei" && p != "kostenlos" && p != "0€" && p != "" {
				continue
			}
		}

		if len(opts.Tags) > 0 {
			match := false
			for _, tag := range opts.Tags {
				for _, et := range e.Tags {
					if strings.EqualFold(tag, et) {
						match = true
					}
				}
			}
			if !match {
				continue
			}
		}

		result = append(result, e)

		if opts.Limit > 0 && len(result) >= opts.Limit {
			break
		}
	}

	return result
}

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/havocked/leipzig-cli/internal/engine"
	"github.com/havocked/leipzig-cli/internal/output"
	"github.com/havocked/leipzig-cli/internal/source/leipzigde"
	"github.com/spf13/cobra"
)

var (
	flagToday    bool
	flagWeekend  bool
	flagCategory string
	flagSearch   string
	flagFree     bool
	flagJSON     bool
	flagFormat   string
	flagLimit    int
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List events in Leipzig",
	RunE:  runEvents,
}

func init() {
	eventsCmd.Flags().BoolVar(&flagToday, "today", false, "Show today's events only")
	eventsCmd.Flags().BoolVar(&flagWeekend, "weekend", false, "Show this weekend's events")
	eventsCmd.Flags().StringVar(&flagCategory, "category", "", "Filter by category (comma-separated)")
	eventsCmd.Flags().StringVar(&flagSearch, "search", "", "Search by name or venue")
	eventsCmd.Flags().BoolVar(&flagFree, "free", false, "Show free events only")
	eventsCmd.Flags().BoolVar(&flagJSON, "json", false, "Output as JSON")
	eventsCmd.Flags().StringVar(&flagFormat, "format", "table", "Output format: table, compact, json")
	eventsCmd.Flags().IntVar(&flagLimit, "limit", 0, "Limit number of results")
	rootCmd.AddCommand(eventsCmd)
}

func runEvents(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	loc, _ := time.LoadLocation("Europe/Berlin")
	now := time.Now().In(loc)

	var from, to time.Time

	if flagToday {
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		to = from.Add(24 * time.Hour)
	} else if flagWeekend {
		// Find next Saturday (or today if already weekend)
		daysUntilSat := (6 - int(now.Weekday()) + 7) % 7
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			daysUntilSat = 0
		}
		sat := time.Date(now.Year(), now.Month(), now.Day()+daysUntilSat, 0, 0, 0, 0, loc)
		from = sat
		to = sat.Add(48 * time.Hour)
	}

	eng := engine.New(leipzigde.New())
	events, err := eng.Fetch(ctx, from, to)
	if err != nil {
		return fmt.Errorf("fetch events: %w", err)
	}

	filtered := engine.Filter(events, engine.FilterOptions{
		Category: flagCategory,
		Search:   flagSearch,
		Free:     flagFree,
		From:     from,
		To:       to,
		Limit:    flagLimit,
	})

	format := flagFormat
	if flagJSON {
		format = "json"
	}

	switch format {
	case "json":
		return output.JSON(os.Stdout, filtered)
	case "compact":
		return output.Compact(os.Stdout, filtered)
	default:
		return output.Table(os.Stdout, filtered)
	}
}

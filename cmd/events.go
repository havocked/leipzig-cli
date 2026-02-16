package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/havocked/leipzig-cli/internal/engine"
	"github.com/havocked/leipzig-cli/internal/output"
	"github.com/havocked/leipzig-cli/internal/source/leipzigde"
	"github.com/havocked/leipzig-cli/internal/source/prinzde"
	"github.com/spf13/cobra"
)

var (
	flagWhen     string
	flagSearch   string
	flagCategory string
	flagAfter    string
	flagJSON     bool
	flagLimit    int
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List events in Leipzig",
	Long: `List events in Leipzig. Defaults to today's events.

Examples:
  leipzig events                          # Today's events
  leipzig events --when weekend           # This weekend
  leipzig events --when tomorrow          # Tomorrow
  leipzig events --search concert         # Search by name/venue
  leipzig events --category family        # Filter by category
  leipzig events --after 16:00            # Events starting at 4 PM or later
  leipzig events --json                   # JSON output for agents
  leipzig events --search jazz --when weekend --json`,
	RunE: runEvents,
}

func init() {
	eventsCmd.Flags().StringVar(&flagWhen, "when", "today", "Time range: today, tomorrow, weekend, week")
	eventsCmd.Flags().StringVarP(&flagSearch, "search", "s", "", "Search by name or venue")
	eventsCmd.Flags().StringVarP(&flagCategory, "category", "c", "", "Filter by category (comma-separated)")
	eventsCmd.Flags().StringVar(&flagAfter, "after", "", "Only events starting at or after this time (HH:MM)")
	eventsCmd.Flags().BoolVar(&flagJSON, "json", false, "Output as JSON")
	eventsCmd.Flags().IntVarP(&flagLimit, "limit", "n", 0, "Limit number of results")
	rootCmd.AddCommand(eventsCmd)
}

func runEvents(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	loc, _ := time.LoadLocation("Europe/Berlin")
	now := time.Now().In(loc)

	from, to := resolveTimeRange(flagWhen, now, loc)

	// Apply --after filter: shift "from" to today/tomorrow at that time
	if flagAfter != "" {
		afterTime, err := time.ParseInLocation("15:04", flagAfter, loc)
		if err != nil {
			return fmt.Errorf("invalid --after time %q (expected HH:MM): %w", flagAfter, err)
		}
		afterFull := time.Date(from.Year(), from.Month(), from.Day(),
			afterTime.Hour(), afterTime.Minute(), 0, 0, loc)
		if afterFull.After(from) {
			from = afterFull
		}
	}

	eng := engine.New(leipzigde.New(), prinzde.New())
	events, err := eng.Fetch(ctx, from, to)
	if err != nil {
		return fmt.Errorf("fetch events: %w", err)
	}

	filtered := engine.Filter(events, engine.FilterOptions{
		Category: flagCategory,
		Search:   flagSearch,
		From:     from,
		To:       to,
		Limit:    flagLimit,
	})

	if flagJSON {
		return output.JSON(os.Stdout, filtered)
	}
	return output.Table(os.Stdout, filtered)
}

func resolveTimeRange(when string, now time.Time, loc *time.Location) (from, to time.Time) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	switch when {
	case "tomorrow":
		from = today.Add(24 * time.Hour)
		to = from.Add(24 * time.Hour)
	case "weekend":
		daysUntilSat := (6 - int(now.Weekday()) + 7) % 7
		if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
			daysUntilSat = 0
		}
		from = today.Add(time.Duration(daysUntilSat) * 24 * time.Hour)
		to = from.Add(48 * time.Hour)
	case "week":
		from = today
		to = today.Add(7 * 24 * time.Hour)
	default: // "today"
		from = today
		to = today.Add(24 * time.Hour)
	}
	return
}

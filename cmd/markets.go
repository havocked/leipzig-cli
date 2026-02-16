package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/havocked/leipzig-cli/internal/market"
	"github.com/spf13/cobra"
)

var (
	marketsDay  string
	marketsJSON bool
)

var marketsCmd = &cobra.Command{
	Use:   "markets",
	Short: "Show Leipzig's weekly markets (Wochenmärkte)",
	Long:  "Display weekly market schedules for Leipzig's 16 Wochenmärkte.",
	RunE:  runMarkets,
}

func init() {
	marketsCmd.Flags().StringVar(&marketsDay, "day", "today", "Filter by day: today, tomorrow, monday-sunday, or all")
	marketsCmd.Flags().BoolVar(&marketsJSON, "json", false, "JSON output")
	rootCmd.AddCommand(marketsCmd)
}

func parseDay(s string) (time.Weekday, bool, error) {
	switch strings.ToLower(s) {
	case "today":
		return time.Now().Weekday(), false, nil
	case "tomorrow":
		return time.Now().Add(24 * time.Hour).Weekday(), false, nil
	case "all":
		return 0, true, nil
	case "monday", "mon":
		return time.Monday, false, nil
	case "tuesday", "tue":
		return time.Tuesday, false, nil
	case "wednesday", "wed":
		return time.Wednesday, false, nil
	case "thursday", "thu":
		return time.Thursday, false, nil
	case "friday", "fri":
		return time.Friday, false, nil
	case "saturday", "sat":
		return time.Saturday, false, nil
	case "sunday", "sun":
		return time.Sunday, false, nil
	default:
		return 0, false, fmt.Errorf("unknown day: %s", s)
	}
}

func dayLabel(s string, day time.Weekday) string {
	switch strings.ToLower(s) {
	case "today":
		return fmt.Sprintf("today (%s)", day)
	case "tomorrow":
		return fmt.Sprintf("tomorrow (%s)", day)
	default:
		return day.String()
	}
}

func runMarkets(cmd *cobra.Command, args []string) error {
	day, all, err := parseDay(marketsDay)
	if err != nil {
		return err
	}

	if all {
		return printAll()
	}

	markets := market.ForDay(day)

	if marketsJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(markets)
	}

	if len(markets) == 0 {
		fmt.Printf("No markets open on %s.\n", dayLabel(marketsDay, day))
		return nil
	}

	fmt.Printf("Markets open %s:\n\n", dayLabel(marketsDay, day))
	for _, m := range markets {
		fmt.Printf("🛍️  %s  %s\n", market.FormatTime(m.Open, m.Close), m.Name)
	}
	fmt.Println("\nUse --day all to see the full weekly schedule.")
	return nil
}

func printAll() error {
	allDays := market.AllByDay()
	order := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

	if marketsJSON {
		out := map[string][]market.MarketDay{}
		for _, d := range order {
			if list, ok := allDays[d]; ok {
				out[d.String()] = list
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Println("Leipzig Weekly Markets (Wochenmärkte):\n")
	for _, d := range order {
		list, ok := allDays[d]
		if !ok {
			continue
		}
		fmt.Printf("📅 %s:\n", d)
		for _, m := range list {
			fmt.Printf("   🛍️  %s  %s\n", market.FormatTime(m.Open, m.Close), m.Name)
		}
		fmt.Println()
	}
	return nil
}

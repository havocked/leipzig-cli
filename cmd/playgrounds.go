package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/havocked/leipzig-cli/internal/playground"
	"github.com/spf13/cobra"
)

var (
	pgDistrict string
	pgSearch   string
	pgJSON     bool
	pgLimit    int
)

var playgroundsCmd = &cobra.Command{
	Use:   "playgrounds",
	Short: "Show Leipzig playgrounds (~320 Spielplätze)",
	Long:  "Scrape and display all public playgrounds from leipzig.de.",
	RunE:  runPlaygrounds,
}

func init() {
	playgroundsCmd.Flags().StringVarP(&pgDistrict, "district", "d", "", "Filter by district or subdistrict (case-insensitive)")
	playgroundsCmd.Flags().StringVarP(&pgSearch, "search", "s", "", "Search by name (case-insensitive)")
	playgroundsCmd.Flags().BoolVar(&pgJSON, "json", false, "JSON output")
	playgroundsCmd.Flags().IntVarP(&pgLimit, "limit", "n", 0, "Max results to show")
	rootCmd.AddCommand(playgroundsCmd)
}

func runPlaygrounds(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(os.Stderr, "Fetching playgrounds from leipzig.de...\n")

	all, err := playground.FetchAll()
	if err != nil {
		return fmt.Errorf("fetch playgrounds: %w", err)
	}

	filtered := all
	if pgDistrict != "" {
		d := strings.ToLower(pgDistrict)
		var out []playground.Playground
		for _, p := range filtered {
			if strings.Contains(strings.ToLower(p.District), d) || strings.Contains(strings.ToLower(p.Subdistrict), d) {
				out = append(out, p)
			}
		}
		filtered = out
	}
	if pgSearch != "" {
		s := strings.ToLower(pgSearch)
		var out []playground.Playground
		for _, p := range filtered {
			if strings.Contains(strings.ToLower(p.Name), s) {
				out = append(out, p)
			}
		}
		filtered = out
	}

	total := len(filtered)
	if pgLimit > 0 && pgLimit < len(filtered) {
		filtered = filtered[:pgLimit]
	}

	if pgJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(filtered)
	}

	if len(filtered) == 0 {
		fmt.Println("No playgrounds found.")
		return nil
	}

	for i, p := range filtered {
		if i > 0 {
			fmt.Println()
		}
		district := p.District
		if p.Subdistrict != "" {
			district += " / " + p.Subdistrict
		}
		fmt.Printf("🛝 %s\n", p.Name)
		if p.Address != "" {
			fmt.Printf("   📍 %s\n", p.Address)
		}
		if district != "" {
			fmt.Printf("   🏘️  %s\n", district)
		}
		if p.MapURL != "" {
			fmt.Printf("   🗺️  %s\n", p.MapURL)
		}
	}

	fmt.Printf("\nShowing %d of %d playgrounds\n", len(filtered), total)
	return nil
}

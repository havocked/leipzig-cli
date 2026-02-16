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
	Short: "Find Leipzig's ~320 public playgrounds",
	Long:  "Search and filter Leipzig's public playgrounds by district or name.",
	RunE:  runPlaygrounds,
}

func init() {
	playgroundsCmd.Flags().StringVarP(&pgDistrict, "district", "d", "", "Filter by district/subdistrict (case-insensitive contains)")
	playgroundsCmd.Flags().StringVarP(&pgSearch, "search", "s", "", "Search by name (case-insensitive contains)")
	playgroundsCmd.Flags().BoolVar(&pgJSON, "json", false, "JSON output")
	playgroundsCmd.Flags().IntVarP(&pgLimit, "limit", "n", 0, "Max results (0=all)")
	rootCmd.AddCommand(playgroundsCmd)
}

func runPlaygrounds(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(os.Stderr, "Fetching playgrounds from leipzig.de...\n")

	all, err := playground.FetchAll()
	if err != nil {
		return fmt.Errorf("fetching playgrounds: %w", err)
	}

	// Filter
	var results []playground.Playground
	for _, p := range all {
		if pgDistrict != "" {
			d := strings.ToLower(pgDistrict)
			if !strings.Contains(strings.ToLower(p.District), d) &&
				!strings.Contains(strings.ToLower(p.Subdistrict), d) {
				continue
			}
		}
		if pgSearch != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(pgSearch)) {
				continue
			}
		}
		results = append(results, p)
	}

	// Limit
	if pgLimit > 0 && len(results) > pgLimit {
		results = results[:pgLimit]
	}

	if pgJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No playgrounds found matching your criteria.")
		return nil
	}

	for _, p := range results {
		fmt.Printf("🛝 %s\n", p.Name)
		if p.Address != "" {
			fmt.Printf("   📍 %s\n", p.Address)
		}
		if p.District != "" {
			loc := p.District
			if p.Subdistrict != "" {
				loc += " / " + p.Subdistrict
			}
			fmt.Printf("   🏘️  %s\n", loc)
		}
		if p.MapURL != "" {
			fmt.Printf("   🗺️  %s\n", p.MapURL)
		}
		fmt.Println()
	}

	fmt.Fprintf(os.Stderr, "Showing %d of %d playgrounds\n", len(results), len(all))
	return nil
}

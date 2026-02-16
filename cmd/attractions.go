package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/havocked/leipzig-cli/internal/attraction"
	"github.com/spf13/cobra"
)

var (
	attrCategory string
	attrSearch   string
	attrJSON     bool
	attrLimit    int
)

var attractionsCmd = &cobra.Command{
	Use:   "attractions",
	Short: "Discover Leipzig's top attractions and sights",
	Long:  "Browse curated Leipzig attractions: landmarks, museums, parks, culture, and more.",
	RunE:  runAttractions,
}

func init() {
	attractionsCmd.Flags().StringVarP(&attrCategory, "category", "c", "", "Filter by category (landmark|museum|church|park|culture|district|family)")
	attractionsCmd.Flags().StringVarP(&attrSearch, "search", "s", "", "Search by name or description")
	attractionsCmd.Flags().BoolVar(&attrJSON, "json", false, "JSON output")
	attractionsCmd.Flags().IntVarP(&attrLimit, "limit", "n", 0, "Max results (0=all)")
	rootCmd.AddCommand(attractionsCmd)
}

var categoryEmoji = map[string]string{
	"landmark": "🏛️",
	"museum":   "🖼️",
	"church":   "⛪",
	"park":     "🌳",
	"culture":  "🎭",
	"district": "🏘️",
	"family":   "👨‍👩‍👧‍👦",
}

func runAttractions(cmd *cobra.Command, args []string) error {
	all := attraction.All()

	var results []attraction.Attraction
	for _, a := range all {
		if attrCategory != "" {
			if !strings.EqualFold(a.Category, attrCategory) {
				continue
			}
		}
		if attrSearch != "" {
			s := strings.ToLower(attrSearch)
			if !strings.Contains(strings.ToLower(a.Name), s) &&
				!strings.Contains(strings.ToLower(a.Description), s) {
				continue
			}
		}
		results = append(results, a)
	}

	if attrLimit > 0 && len(results) > attrLimit {
		results = results[:attrLimit]
	}

	if attrJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No attractions found matching your criteria.")
		return nil
	}

	for _, a := range results {
		emoji := categoryEmoji[a.Category]
		if emoji == "" {
			emoji = "📍"
		}
		fmt.Printf("%s %s\n", emoji, a.Name)
		fmt.Printf("   %s\n", a.Description)
		if a.Address != "" {
			fmt.Printf("   📍 %s\n", a.Address)
		}
		if a.URL != "" {
			fmt.Printf("   🔗 %s\n", a.URL)
		}
		fmt.Println()
	}

	fmt.Fprintf(os.Stderr, "Showing %d of %d attractions\n", len(results), len(all))
	return nil
}

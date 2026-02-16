package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/havocked/leipzig-cli/internal/news"
	"github.com/spf13/cobra"
)

var (
	newsCategory string
	newsSearch   string
	newsPages    int
	newsJSON     bool
	newsLimit    int
)

var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Show recent Leipzig city news",
	Long:  "Fetch and display recent news from leipzig.de/newsarchiv.",
	RunE:  runNews,
}

func init() {
	newsCmd.Flags().StringVarP(&newsCategory, "category", "c", "", "Filter by category (building|services|leisure|health|international|family|children|culture|safety|traffic|economy)")
	newsCmd.Flags().StringVarP(&newsSearch, "search", "s", "", "Search by keyword")
	newsCmd.Flags().IntVar(&newsPages, "pages", 1, "Number of pages to fetch (25 articles/page)")
	newsCmd.Flags().BoolVar(&newsJSON, "json", false, "JSON output")
	newsCmd.Flags().IntVarP(&newsLimit, "limit", "n", 0, "Max results (0=all)")
	rootCmd.AddCommand(newsCmd)
}

func runNews(cmd *cobra.Command, args []string) error {
	opts := news.FetchOptions{
		Pages:  newsPages,
		Search: newsSearch,
	}

	if newsCategory != "" {
		cat := strings.ToLower(newsCategory)
		apiVal, ok := news.Categories[cat]
		if !ok {
			// List valid categories
			var valid []string
			for k := range news.Categories {
				valid = append(valid, k)
			}
			sort.Strings(valid)
			return fmt.Errorf("unknown category %q. Valid: %s", newsCategory, strings.Join(valid, ", "))
		}
		opts.Category = apiVal
	}

	fmt.Fprintf(os.Stderr, "Fetching news from leipzig.de...\n")

	articles, err := news.Fetch(opts)
	if err != nil {
		return fmt.Errorf("fetching news: %w", err)
	}

	if newsLimit > 0 && len(articles) > newsLimit {
		articles = articles[:newsLimit]
	}

	if newsJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(articles)
	}

	if len(articles) == 0 {
		fmt.Println("No news articles found.")
		return nil
	}

	for _, a := range articles {
		fmt.Printf("📰 %s — %s\n", a.Date, a.Title)
		fmt.Printf("   🔗 %s\n\n", a.URL)
	}

	fmt.Fprintf(os.Stderr, "Showing %d articles\n", len(articles))
	return nil
}

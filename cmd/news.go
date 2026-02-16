package cmd

import (
	"encoding/json"
	"fmt"
	"os"
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
	Short: "Show recent news from Leipzig",
	Long: `Scrape and display recent news from leipzig.de/newsarchiv.

Examples:
  leipzig news                  # 25 most recent articles
  leipzig news -n 5             # limit to 5
  leipzig news -c traffic       # filter by category
  leipzig news -s "Spielplatz"  # search by keyword
  leipzig news --json -n 3      # JSON output

Categories: building, services, leisure, health, international,
            family, children, culture, safety, traffic, economy`,
	RunE: runNews,
}

func init() {
	newsCmd.Flags().StringVarP(&newsCategory, "category", "c", "", "Filter by category")
	newsCmd.Flags().StringVarP(&newsSearch, "search", "s", "", "Search by keyword")
	newsCmd.Flags().IntVar(&newsPages, "pages", 1, "Number of pages to fetch (25 per page)")
	newsCmd.Flags().BoolVar(&newsJSON, "json", false, "JSON output")
	newsCmd.Flags().IntVarP(&newsLimit, "limit", "n", 0, "Max number of results")
	rootCmd.AddCommand(newsCmd)
}

func runNews(cmd *cobra.Command, args []string) error {
	articles, err := news.Fetch(news.FetchOptions{
		Pages:    newsPages,
		Category: newsCategory,
		Search:   newsSearch,
	})
	if err != nil {
		return err
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

	label := "articles"
	if newsCategory != "" {
		label = fmt.Sprintf("articles [%s]", strings.ToLower(newsCategory))
	}
	fmt.Printf("Showing %d %s\n", len(articles), label)
	return nil
}

package news

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.leipzig.de"
const archiveURL = baseURL + "/newsarchiv"

// CategoryMap maps user-friendly English names to API category values.
var CategoryMap = map[string]string{
	"building":      "128",
	"services":      "24",
	"leisure":       "49",
	"health":        "28",
	"international": "244",
	"family":        "26",
	"children":      "518",
	"culture":       "143",
	"safety":        "52",
	"traffic":       "51",
	"economy":       "127",
}

// FetchOptions configures the news fetch.
type FetchOptions struct {
	Pages    int
	Category string // user-friendly name (e.g. "traffic")
	Search   string
}

// Fetch scrapes news articles from leipzig.de/newsarchiv.
func Fetch(opts FetchOptions) ([]Article, error) {
	if opts.Pages < 1 {
		opts.Pages = 1
	}

	client := &http.Client{Timeout: 30 * time.Second}
	var articles []Article

	for page := 0; page < opts.Pages; page++ {
		if page > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		u, err := buildURL(opts, page)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "leipzig-cli/1.0")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch page %d: %w", page+1, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("page %d: HTTP %d", page+1, resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("parse page %d: %w", page+1, err)
		}

		doc.Find("div.card.news-result-card").Each(func(_ int, s *goquery.Selection) {
			a := Article{}

			a.Date = strings.TrimSpace(s.Find("div.card-text").First().Text())

			link := s.Find("h3.card-title a")
			a.Title = strings.TrimSpace(strings.ReplaceAll(link.Find("span").Text(), "arrow_forward", ""))
			if href, ok := link.Attr("href"); ok {
				if strings.HasPrefix(href, "/") {
					a.URL = baseURL + href
				} else {
					a.URL = href
				}
			}

			if src, ok := s.Find("img").Attr("src"); ok {
				if strings.HasPrefix(src, "/") {
					a.ImageURL = baseURL + src
				} else {
					a.ImageURL = src
				}
			}

			if a.Title != "" {
				articles = append(articles, a)
			}
		})
	}

	return articles, nil
}

func buildURL(opts FetchOptions, page int) (string, error) {
	params := url.Values{}

	if page > 0 {
		params.Set("mksearch[pb-search327975-pointer]", fmt.Sprintf("%d", page))
	}

	if opts.Category != "" {
		catVal, ok := CategoryMap[strings.ToLower(opts.Category)]
		if !ok {
			return "", fmt.Errorf("unknown category %q (available: building, services, leisure, health, international, family, children, culture, safety, traffic, economy)", opts.Category)
		}
		params.Set("mksearch[category][]", catVal)
		params.Set("mksearch[submitted]", "1")
	}

	if opts.Search != "" {
		params.Set("mksearch[term]", opts.Search)
		params.Set("mksearch[submitted]", "1")
	}

	if len(params) == 0 {
		return archiveURL, nil
	}
	return archiveURL + "?" + params.Encode(), nil
}

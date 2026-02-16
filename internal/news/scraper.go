package news

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL   = "https://www.leipzig.de/newsarchiv"
	userAgent = "leipzig-cli/1.0 (+https://github.com/havocked/leipzig-cli)"
	delay     = 500 * time.Millisecond
)

type FetchOptions struct {
	Pages    int
	Category string // API value (e.g. "51")
	Search   string
}

func Fetch(opts FetchOptions) ([]Article, error) {
	if opts.Pages <= 0 {
		opts.Pages = 1
	}

	var articles []Article
	for page := 0; page < opts.Pages; page++ {
		if page > 0 {
			time.Sleep(delay)
		}
		pageArticles, err := fetchPage(page, opts)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", page+1, err)
		}
		articles = append(articles, pageArticles...)
		if len(pageArticles) == 0 {
			break
		}
	}

	return articles, nil
}

func fetchPage(page int, opts FetchOptions) ([]Article, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("mksearch[pb-search327975-pointer]", fmt.Sprintf("%d", page))
	}
	if opts.Search != "" {
		params.Set("mksearch[term]", opts.Search)
		params.Set("mksearch[submitted]", "1")
	}
	if opts.Category != "" {
		params.Set("mksearch[category][]", opts.Category)
		params.Set("mksearch[submitted]", "1")
	}

	pageURL := baseURL
	if encoded := params.Encode(); encoded != "" {
		pageURL += "?" + encoded
	}

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var articles []Article

	doc.Find(".card.news-result-card").Each(func(_ int, s *goquery.Selection) {
		a := parseArticle(s)
		if a.Title != "" {
			articles = append(articles, a)
		}
	})

	return articles, nil
}

func parseArticle(s *goquery.Selection) Article {
	var a Article

	// Date
	a.Date = strings.TrimSpace(s.Find(".card-text").First().Text())

	// Title and URL
	link := s.Find("h3.card-title a")
	a.Title = strings.TrimSpace(link.Find("span").First().Text())
	if href, ok := link.Attr("href"); ok {
		if strings.HasPrefix(href, "/") {
			a.URL = "https://www.leipzig.de" + href
		} else {
			a.URL = href
		}
	}

	// Image (optional)
	if src, ok := s.Find("img").Attr("src"); ok {
		a.ImageURL = src
	}

	return a
}

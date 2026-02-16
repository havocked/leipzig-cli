package playground

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL   = "https://www.leipzig.de/kultur-und-freizeit/spielplaetze"
	perPage   = 25
	userAgent = "leipzig-cli/1.0"
)

func pageURL(page int) string {
	return fmt.Sprintf("%s?tx_lepurpose%%5Bfilter%%5D%%5Bpage%%5D=%d&tx_lepurpose%%5Bfiltered%%5D=1", baseURL, page)
}

var totalRe = regexp.MustCompile(`von\s+<strong>(\d+)</strong>`)

func FetchAll() ([]Playground, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	playgrounds, total, err := fetchPage(client, 1)
	if err != nil {
		return nil, fmt.Errorf("page 1: %w", err)
	}

	if total == 0 {
		total = len(playgrounds)
	}

	totalPages := (total + perPage - 1) / perPage
	fmt.Fprintf(nil, "") // removed debug

	for p := 2; p <= totalPages; p++ {
		time.Sleep(500 * time.Millisecond)
		pg, _, err := fetchPage(client, p)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", p, err)
		}
		playgrounds = append(playgrounds, pg...)
	}

	sort.Slice(playgrounds, func(i, j int) bool {
		return playgrounds[i].Name < playgrounds[j].Name
	})

	return playgrounds, nil
}

func fetchPage(client *http.Client, page int) ([]Playground, int, error) {
	req, err := http.NewRequest("GET", pageURL(page), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "de-DE,de;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	// Parse total: "von <strong>320</strong> Ergebnissen"
	total := 0
	doc.Find(".result-info").Each(func(_ int, s *goquery.Selection) {
		html, _ := s.Html()
		matches := totalRe.FindStringSubmatch(html)
		if len(matches) > 1 {
			total, _ = strconv.Atoi(matches[1])
		}
	})

	var playgrounds []Playground

	doc.Find("article.filterlist-teaser-card").Each(func(_ int, card *goquery.Selection) {
		pg := Playground{}

		// Name & URL
		link := card.Find("a.filterlist-detailpage")
		pg.Name = strings.TrimSpace(link.Find("span").First().Text())
		if href, ok := link.Attr("href"); ok {
			if strings.HasPrefix(href, "/") {
				pg.DetailURL = "https://www.leipzig.de" + href
			} else {
				pg.DetailURL = href
			}
		}

		// Address & District from the span with class=""
		card.Find("div.d-flex span").Each(func(_ int, s *goquery.Selection) {
			cls, exists := s.Attr("class")
			if !exists || cls != "" {
				return
			}
			if pg.Address != "" {
				return
			}

			// Get direct text nodes (before <br/> and <ul>)
			// The span contains: "Address <br/> <ul>..."
			// Use goquery to get the text content, then parse
			fullText := s.Contents().FilterFunction(func(_ int, sel *goquery.Selection) bool {
				return goquery.NodeName(sel) == "#text"
			}).First().Text()
			pg.Address = strings.TrimSpace(fullText)

			// District from nested ul li
			districtText := strings.TrimSpace(s.Find("ul.list-unstyled li").First().Text())
			if strings.Contains(districtText, " / ") {
				parts := strings.SplitN(districtText, " / ", 2)
				pg.District = strings.TrimSpace(parts[0])
				pg.Subdistrict = strings.TrimSpace(parts[1])
			} else {
				pg.District = districtText
			}
		})

		if pg.Address != "" {
			pg.MapURL = pg.BuildMapURL()
		}

		if pg.Name != "" {
			playgrounds = append(playgrounds, pg)
		}
	})

	return playgrounds, total, nil
}

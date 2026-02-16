package playground

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL   = "https://www.leipzig.de/kultur-und-freizeit/spielplaetze"
	userAgent = "leipzig-cli/1.0 (+https://github.com/havocked/leipzig-cli)"
	pageSize  = 25
	delay     = 500 * time.Millisecond
)

func FetchAll() ([]Playground, error) {
	// Fetch page 1 to get total count
	playgrounds, total, err := fetchPage(1)
	if err != nil {
		return nil, fmt.Errorf("fetching page 1: %w", err)
	}

	if total <= pageSize {
		sort.Slice(playgrounds, func(i, j int) bool {
			return playgrounds[i].Name < playgrounds[j].Name
		})
		return playgrounds, nil
	}

	totalPages := (total + pageSize - 1) / pageSize

	for page := 2; page <= totalPages; page++ {
		time.Sleep(delay)
		pg, _, err := fetchPage(page)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", page, err)
		}
		playgrounds = append(playgrounds, pg...)
		if len(pg) == 0 {
			break
		}
	}

	sort.Slice(playgrounds, func(i, j int) bool {
		return playgrounds[i].Name < playgrounds[j].Name
	})
	return playgrounds, nil
}

func fetchPage(page int) ([]Playground, int, error) {
	pageURL := fmt.Sprintf("%s?tx_lepurpose[filtered]=1&tx_lepurpose[filter][page]=%d#c313548", baseURL, page)

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
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

	var playgrounds []Playground

	doc.Find("article.filterlist-teaser-card").Each(func(_ int, s *goquery.Selection) {
		p := parsePlayground(s)
		if p.Name != "" {
			playgrounds = append(playgrounds, p)
		}
	})

	// Parse total count from "1 - 25 von 320 Ergebnissen"
	total := 0
	doc.Find(".result-info strong").Each(func(i int, s *goquery.Selection) {
		if i == 1 { // second <strong> is total
			fmt.Sscanf(s.Text(), "%d", &total)
		}
	})

	return playgrounds, total, nil
}

func parsePlayground(s *goquery.Selection) Playground {
	var p Playground

	// Name from link text
	link := s.Find("a.filterlist-detailpage")
	p.Name = strings.TrimSpace(link.Find("span").First().Text())

	// Detail URL
	if href, ok := link.Attr("href"); ok {
		p.DetailURL = "https://www.leipzig.de" + href
	}

	// Address - text before <br/> in the span
	addrSpan := s.Find(".card-body .d-flex span").First()
	// Get the raw HTML and split on <br/>
	addrHTML, _ := addrSpan.Html()
	parts := strings.SplitN(addrHTML, "<br/>", 2)
	if len(parts) > 0 {
		p.Address = strings.TrimSpace(stripTags(parts[0]))
	}

	// District from list item
	distText := strings.TrimSpace(s.Find(".card-body ul.list-unstyled li").First().Text())
	distParts := strings.SplitN(distText, " / ", 2)
	if len(distParts) == 2 {
		p.District = strings.TrimSpace(distParts[0])
		p.Subdistrict = strings.TrimSpace(distParts[1])
	} else if len(distParts) == 1 {
		p.District = strings.TrimSpace(distParts[0])
	}

	if p.Address != "" {
		p.MapURL = MakeMapURL(p.Address)
	}

	return p
}

func stripTags(s string) string {
	// Simple tag stripper for inline HTML
	result := strings.Builder{}
	inTag := false
	for _, c := range s {
		if c == '<' {
			inTag = true
			continue
		}
		if c == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(c)
		}
	}
	return result.String()
}

package playground

import (
	"fmt"
	"net/url"
	"strings"
)

type Playground struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
	DetailURL   string `json:"detail_url"`
	MapURL      string `json:"map_url"`
}

func MakeMapURL(address string) string {
	q := strings.TrimSpace(address) + ", Leipzig"
	return fmt.Sprintf("https://maps.google.com/?q=%s", url.QueryEscape(q))
}

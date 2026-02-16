package news

type Article struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	URL      string `json:"url"`
	Category string `json:"category,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// Category mapping: user-friendly name → API value
var Categories = map[string]string{
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

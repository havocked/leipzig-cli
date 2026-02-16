package news

// Article represents a single news article from leipzig.de/newsarchiv.
type Article struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	URL      string `json:"url"`
	Category string `json:"category,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

package model

import "strings"

// nameKeywords maps lowercase keywords found in event names/venues to categories.
// Order matters: first match wins. More specific terms come first.
var nameKeywords = []struct {
	keywords []string
	category string
}{
	// Family / Kids
	{[]string{"kinder", "familien", "winterferien", "mädchenwerkstatt", "jugend"}, CategoryFamily},

	// Sport / Fitness
	{[]string{"gymnastik", "nordic walking", "fitness", "yoga", "training", "sport",
		"lauf", "running", "schwimm", "eisbad", "turnen", "bewegung"}, CategorySport},

	// Concert / Music
	{[]string{"konzert", "chor", "blockflöte", "musikalisch", "singen",
		"volkslied", "band", "orchester", "live music", "livemusik", "dj"}, CategoryConcert},

	// Theater / Comedy / Cabaret
	{[]string{"theater", "kabarett", "comedy", "varieté", "fasching",
		"karneval", "rosenmontag", "bühne", "schauspiel", "impro"}, CategoryTheater},

	// Exhibition
	{[]string{"ausstellung", "exhibition", "orchidee", "museum", "galerie",
		"panorama", "schau"}, CategoryExhibition},

	// Market
	{[]string{"markt", "flohmarkt", "weihnachtsmarkt", "messe"}, CategoryMarket},

	// Culture (broad: tours, talks, language, creative, libraries)
	{[]string{"führung", "rundgang", "lesung", "vortrag", "bibliothek",
		"englisch", "konversation", "sprachkurs", "stammtisch",
		"mal", "zeichn", "kreativ", "bastel", "handarbeit",
		"natur-erlebnis", "werkstatt", "fahrrad"}, CategoryCulture},

	// Nightlife
	{[]string{"club", "party", "nachtleben", "disco", "rave"}, CategoryNightlife},
}

// InferCategory tries to guess a category from the event name and venue
// when the source didn't provide one (or mapped to "other").
// Returns the inferred category, or CategoryOther if nothing matches.
func InferCategory(name, venue string) string {
	combined := strings.ToLower(name + " " + venue)

	for _, rule := range nameKeywords {
		for _, kw := range rule.keywords {
			if strings.Contains(combined, kw) {
				return rule.category
			}
		}
	}

	return CategoryOther
}

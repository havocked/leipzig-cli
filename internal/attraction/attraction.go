package attraction

import (
	"fmt"
	"net/url"
)

type Attraction struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Address     string `json:"address,omitempty"`
	URL         string `json:"url,omitempty"`
	MapURL      string `json:"map_url"`
}

func makeMapURL(address string) string {
	return fmt.Sprintf("https://maps.google.com/?q=%s", url.QueryEscape(address))
}

// All returns all curated Leipzig attractions
func All() []Attraction {
	return attractions
}

var attractions = []Attraction{
	// Landmarks
	{Name: "Völkerschlachtdenkmal", Category: "landmark", Description: "Monument to the 1813 Battle of the Nations, panoramic viewing platform", Address: "Str. des 18. Oktober 100, 04299 Leipzig", URL: "https://www.stadtgeschichtliches-museum-leipzig.de/voelkerschlachtdenkmal", MapURL: makeMapURL("Völkerschlachtdenkmal Leipzig")},
	{Name: "Neues Rathaus", Category: "landmark", Description: "One of Germany's largest city halls, tower with city views", Address: "Martin-Luther-Ring 4-6, 04109 Leipzig", MapURL: makeMapURL("Neues Rathaus Leipzig")},
	{Name: "Altes Rathaus", Category: "landmark", Description: "Renaissance building on the Markt, houses city history museum", Address: "Markt 1, 04109 Leipzig", MapURL: makeMapURL("Altes Rathaus Leipzig")},
	{Name: "Augustusplatz", Category: "landmark", Description: "One of Europe's largest city squares, home to opera, university, and Gewandhaus", Address: "Augustusplatz, 04109 Leipzig", MapURL: makeMapURL("Augustusplatz Leipzig")},
	{Name: "Mädler-Passage", Category: "landmark", Description: "Historic shopping arcade with Auerbachs Keller (Goethe's Faust)", Address: "Grimmaische Str. 2-4, 04109 Leipzig", MapURL: makeMapURL("Mädler-Passage Leipzig")},

	// Museums
	{Name: "Zeitgeschichtliches Forum", Category: "museum", Description: "Free museum on German history from 1945 to present", Address: "Grimmaische Str. 6, 04109 Leipzig", URL: "https://www.hdg.de/zeitgeschichtliches-forum", MapURL: makeMapURL("Zeitgeschichtliches Forum Leipzig")},
	{Name: "Grassi Museum", Category: "museum", Description: "Three museums: Applied Arts, Musical Instruments, Ethnography", Address: "Johannisplatz 5-11, 04103 Leipzig", URL: "https://www.grassimuseum.de", MapURL: makeMapURL("Grassi Museum Leipzig")},
	{Name: "Museum der bildenden Künste", Category: "museum", Description: "Fine arts museum with works from medieval to contemporary", Address: "Katharinenstr. 10, 04109 Leipzig", URL: "https://mdbk.de", MapURL: makeMapURL("Museum der bildenden Künste Leipzig")},
	{Name: "Bach-Museum", Category: "museum", Description: "Museum dedicated to Johann Sebastian Bach's life in Leipzig", Address: "Thomaskirchhof 15/16, 04109 Leipzig", URL: "https://www.bachmuseumleipzig.de", MapURL: makeMapURL("Bach-Museum Leipzig")},
	{Name: "Panometer", Category: "museum", Description: "360° panoramic art installations in a former gasometer", Address: "Richard-Lehmann-Str. 114, 04275 Leipzig", URL: "https://www.panometer.de", MapURL: makeMapURL("Panometer Leipzig")},
	{Name: "Zeitgenössische Kunst (GfZK)", Category: "museum", Description: "Contemporary art gallery with rotating exhibitions", Address: "Karl-Tauchnitz-Str. 9-11, 04107 Leipzig", URL: "https://www.gfzk.de", MapURL: makeMapURL("GfZK Leipzig")},
	{Name: "Naturkundemuseum", Category: "museum", Description: "Natural history museum with regional focus", Address: "Lortzingstr. 3, 04105 Leipzig", URL: "https://www.naturkundemuseum.leipzig.de", MapURL: makeMapURL("Naturkundemuseum Leipzig")},

	// Churches
	{Name: "Thomaskirche", Category: "church", Description: "Bach's church, home of the famous Thomanerchor (boys' choir)", Address: "Thomaskirchhof 18, 04109 Leipzig", URL: "https://www.thomaskirche.org", MapURL: makeMapURL("Thomaskirche Leipzig")},
	{Name: "Nikolaikirche", Category: "church", Description: "Key site of the 1989 Peaceful Revolution, Monday demonstrations", Address: "Nikolaikirchhof 3, 04109 Leipzig", URL: "https://www.nikolaikirche.de", MapURL: makeMapURL("Nikolaikirche Leipzig")},

	// Parks & Nature
	{Name: "Zoo Leipzig", Category: "park", Description: "One of Europe's best zoos with Gondwanaland tropical hall", Address: "Pfaffendorfer Str. 29, 04105 Leipzig", URL: "https://www.zoo-leipzig.de", MapURL: makeMapURL("Zoo Leipzig")},
	{Name: "Clara-Zetkin-Park", Category: "park", Description: "Large city park, popular for jogging, picnics, and events", Address: "Clara-Zetkin-Park, 04107 Leipzig", MapURL: makeMapURL("Clara-Zetkin-Park Leipzig")},
	{Name: "Leipziger Auwald", Category: "park", Description: "One of the largest urban floodplain forests in Europe", Address: "Auwald, Leipzig", MapURL: makeMapURL("Auwald Leipzig")},
	{Name: "Cospudener See", Category: "park", Description: "Recreational lake south of Leipzig, beaches and water sports", Address: "Cospudener See, Leipzig", MapURL: makeMapURL("Cospudener See Leipzig")},
	{Name: "Fockeberg", Category: "park", Description: "Hill from WWII rubble with panoramic city views", Address: "Fockestraße, 04275 Leipzig", MapURL: makeMapURL("Fockeberg Leipzig")},
	{Name: "Botanischer Garten", Category: "park", Description: "University botanical garden with greenhouses and rare plants", Address: "Linnéstr. 1, 04103 Leipzig", URL: "https://www.uni-leipzig.de/botanischer-garten", MapURL: makeMapURL("Botanischer Garten Leipzig")},

	// Culture & Entertainment
	{Name: "Gewandhaus", Category: "culture", Description: "World-renowned concert hall, home of the Gewandhausorchester", Address: "Augustusplatz 8, 04109 Leipzig", URL: "https://www.gewandhausorchester.de", MapURL: makeMapURL("Gewandhaus Leipzig")},
	{Name: "Oper Leipzig", Category: "culture", Description: "One of Germany's oldest opera houses", Address: "Augustusplatz 12, 04109 Leipzig", URL: "https://www.oper-leipzig.de", MapURL: makeMapURL("Oper Leipzig")},
	{Name: "Spinnerei", Category: "culture", Description: "Former cotton mill, now a massive art gallery and studio complex", Address: "Spinnereistr. 7, 04179 Leipzig", URL: "https://www.spinnerei.de", MapURL: makeMapURL("Spinnerei Leipzig")},
	{Name: "Schaubühne Lindenfels", Category: "culture", Description: "Independent cinema and cultural venue", Address: "Karl-Heine-Str. 50, 04229 Leipzig", URL: "https://www.schaubuehne.com", MapURL: makeMapURL("Schaubühne Lindenfels Leipzig")},

	// Neighborhoods & Districts
	{Name: "Karl-Liebknecht-Straße (KarLi)", Category: "district", Description: "Leipzig's most vibrant bar and restaurant street", Address: "Karl-Liebknecht-Str., 04275 Leipzig", MapURL: makeMapURL("Karl-Liebknecht-Straße Leipzig")},
	{Name: "Plagwitz / Karl-Heine-Kanal", Category: "district", Description: "Trendy district with canal walks, galleries, and cafés", Address: "Karl-Heine-Kanal, 04229 Leipzig", MapURL: makeMapURL("Karl-Heine-Kanal Leipzig")},
	{Name: "Leipziger Baumwollspinnerei", Category: "district", Description: "Art district in the former cotton mill complex", Address: "Spinnereistr. 7, 04179 Leipzig", MapURL: makeMapURL("Baumwollspinnerei Leipzig")},

	// Family
	{Name: "Belantis", Category: "family", Description: "Largest amusement park in eastern Germany", Address: "Zur Weißen Mark 1, 04249 Leipzig", URL: "https://www.belantis.de", MapURL: makeMapURL("Belantis Leipzig")},
	{Name: "Wildpark Leipzig", Category: "family", Description: "Free wildlife park in the Auwald with deer, wild boar, and birds", Address: "Koburger Str. 12, 04277 Leipzig", MapURL: makeMapURL("Wildpark Leipzig")},
}

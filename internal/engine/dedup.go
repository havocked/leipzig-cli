package engine

import (
	"strings"

	"github.com/havocked/leipzig-cli/internal/model"
)

// Dedup removes cross-source duplicate events. It groups by normalized venue +
// date + hour, then checks name similarity within each group.
func Dedup(events []model.Event) []model.Event {
	type key struct {
		venue string
		date  string // "2006-01-02"
		hour  int    // -1 means midnight/zero-time
	}

	groups := make(map[key][]int) // key -> indices
	for i, e := range events {
		k := key{
			venue: normalizeVenue(e.Venue),
			date:  e.StartTime.Format("2006-01-02"),
			hour:  e.StartTime.Hour(),
		}
		groups[k] = append(groups[k], i)
	}

	removed := make(map[int]bool)

	for _, indices := range groups {
		if len(indices) < 2 {
			continue
		}
		for i := 0; i < len(indices); i++ {
			if removed[indices[i]] {
				continue
			}
			for j := i + 1; j < len(indices); j++ {
				if removed[indices[j]] {
					continue
				}
				a, b := events[indices[i]], events[indices[j]]
				// Only dedup cross-source
				if a.Source == b.Source {
					continue
				}
				if namesMatch(a.Name, b.Name) {
					loser := pickLoser(a, b, indices[i], indices[j])
					removed[loser] = true
				}
			}
		}
	}

	result := make([]model.Event, 0, len(events)-len(removed))
	for i, e := range events {
		if !removed[i] {
			result = append(result, e)
		}
	}
	return result
}

func normalizeVenue(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.TrimSuffix(v, "leipzig")
	v = strings.TrimSuffix(v, ",")
	v = strings.TrimSpace(v)
	return v
}

func normalizeName(n string) string {
	n = strings.ToLower(strings.TrimSpace(n))
	n = strings.ReplaceAll(n, "\u2013", "-") // en-dash
	n = strings.ReplaceAll(n, "\u2014", "-") // em-dash
	// collapse whitespace
	fields := strings.Fields(n)
	return strings.Join(fields, " ")
}

func namesMatch(a, b string) bool {
	na, nb := normalizeName(a), normalizeName(b)
	if na == nb {
		return true
	}
	if strings.Contains(na, nb) || strings.Contains(nb, na) {
		return true
	}
	// Long common prefix (20+ chars)
	const prefixLen = 20
	if len(na) >= prefixLen && len(nb) >= prefixLen && na[:prefixLen] == nb[:prefixLen] {
		return true
	}
	// Shared distinctive words — catches renamed same-event cases
	// Require 2+ shared significant words, or 1 shared word that's 10+ chars
	wordsA := significantWords(na)
	wordsB := significantWords(nb)
	shared := 0
	longShared := false
	for w := range wordsA {
		if wordsB[w] {
			shared++
			if len(w) >= 10 {
				longShared = true
			}
		}
	}
	if shared >= 2 || longShared {
		return true
	}
	return false
}

// pickLoser returns the index of the event to remove.
// Prefer keeping leipzig.de; if tied on source, keep the richer one.
func pickLoser(a, b model.Event, idxA, idxB int) int {
	scoreA := richness(a)
	scoreB := richness(b)

	// Prefer leipzig.de as canonical source
	if a.Source == "leipzig.de" && b.Source != "leipzig.de" {
		scoreA += 100
	}
	if b.Source == "leipzig.de" && a.Source != "leipzig.de" {
		scoreB += 100
	}

	if scoreA >= scoreB {
		return idxB // remove B
	}
	return idxA // remove A
}

// significantWords returns distinctive words (6+ chars, not common German stopwords).
func significantWords(s string) map[string]bool {
	stop := map[string]bool{
		"leipzig": true, "leipziger": true, "veranstaltung": true,
		"durch": true, "einen": true, "einer": true, "diesem": true,
		"dieser": true, "werden": true, "können": true,
	}
	result := make(map[string]bool)
	for _, w := range strings.Fields(s) {
		if len(w) >= 6 && !stop[w] {
			result[w] = true
		}
	}
	return result
}

func richness(e model.Event) int {
	score := len(e.Description)
	if e.ImageURL != "" {
		score += 50
	}
	if !e.EndTime.IsZero() {
		score += 20
	}
	if e.Price != "" {
		score += 10
	}
	return score
}

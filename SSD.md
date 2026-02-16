# SSD — leipzig-cli

## What Is This?
A CLI tool for discovering events and activities in Leipzig. Concerts, family outings, flea markets, exhibitions, nightlife — everything happening in the city, from multiple sources, in one unified format.

**Primary use case:** "What should we do this weekend?" — asked by a human or an AI agent.

## Why?
- Leipzig has great events scattered across many websites with no single aggregator
- KOKA36 (the local concert listing site) has been broken for months
- We want one command to search across all sources with filtering by date, category, family-friendliness, etc.
- Agent-friendly: Ori can use this every Friday to suggest weekend plans

## Architecture

### Core Principle: Canonical Event Model
Every source (website, API) speaks a different language. Each source gets an **adapter** that translates its raw data into our canonical `Event` model. Everything downstream (filtering, sorting, formatting) only knows about `Event` — never about source-specific structures.

```
Sources (adapters)        Canonical Model         Output
┌──────────────┐         ┌─────────────┐         ┌──────────────┐
│ leipzig.de   │──adapt─▶│             │         │ filter       │
│ leipzig-im   │──adapt─▶│   Event[]   │────────▶│ sort         │
│ songkick     │──adapt─▶│             │         │ format (JSON/│
│ (future...)  │──adapt─▶│             │         │  table/text) │
└──────────────┘         └─────────────┘         └──────────────┘
```

This is the **Anti-Corruption Layer** pattern — each adapter isolates the messy external world from our clean domain model.

### Canonical Event Model

```go
type Event struct {
    Name        string    // Event title
    Description string    // Short description (optional, may be empty)
    StartTime   time.Time // When it starts
    EndTime     time.Time // When it ends (zero value = unknown/single-day)
    Venue       string    // Venue name (e.g., "Werk 2", "Panometer")
    Address     string    // Street address (optional)
    Category    string    // Canonical category (see below)
    Tags        []string  // Flexible labels: "outdoor", "kid-friendly", "english"
    Price       string    // "free", "12€", "unknown"
    URL         string    // Link to event details/tickets
    Source      string    // Provider ID: "leipzig.de", "leipzig-im", "songkick"
}
```

### Canonical Categories
Fixed enum, mapped per adapter. Evolves over time based on real data:

| Category     | What it covers |
|-------------|---------------|
| `concert`    | Live music, gigs, DJ sets |
| `theater`    | Theater, comedy, cabaret, performance |
| `exhibition` | Art exhibitions, museum shows, galleries |
| `family`     | Kid-friendly events, children's theater, family activities |
| `market`     | Flea markets, food markets, craft markets |
| `sport`      | Sports events, runs, outdoor activities |
| `culture`    | Readings, lectures, festivals, film screenings |
| `nightlife`  | Club nights, parties, bar events |
| `other`      | Anything that doesn't fit above |

Each adapter maintains its own category mapping table. Unknown source categories fall back to `other`.

### Source Adapter Interface

```go
type Source interface {
    // ID returns the source identifier (e.g., "leipzig.de")
    ID() string

    // Fetch retrieves events within the given time range
    Fetch(ctx context.Context, from, to time.Time) ([]Event, error)
}
```

All sources implement this interface. The engine calls `Fetch` on each enabled source, merges results, deduplicates, and passes to filtering/output.

## Commands

```bash
# List events (default: today + next 7 days)
leipzig events

# This weekend
leipzig events --weekend

# Filter by category
leipzig events --category concert
leipzig events --category family,market

# Filter by date range
leipzig events --from 2026-02-20 --to 2026-02-22

# Today only
leipzig events --today

# Search by text (name/description/venue)
leipzig events --search "Werk 2"
leipzig events --search "flohmarkt"

# Filter by tags
leipzig events --tag outdoor
leipzig events --tag kid-friendly

# Free events only
leipzig events --free

# Limit results
leipzig events --limit 10

# Output formats
leipzig events --json                 # JSON array (agent-friendly)
leipzig events --format table         # Human-readable table (default)
leipzig events --format compact       # One-liner per event

# Source management
leipzig sources                       # List available sources and status
leipzig sources --enable songkick
leipzig sources --disable leipzig-im

# Cache management
leipzig cache clear
leipzig cache status
```

## Output

### Default (table)
```
Sat 22 Feb  18:00  concert    Jinjer – European Duél Tour        Felsenkeller     12€
Sat 22 Feb  10:00  market     Flohmarkt Plagwitz                 Markthalle        free
Sat 22 Feb  14:00  family     Familienführung Antarktis          Panometer         8€
Sun 23 Feb  11:00  culture    Leipziger Buchmesse Preview        Neues Rathaus     free
```

### JSON (for agents / piping)
```json
[
  {
    "name": "Jinjer – European Duél Tour",
    "startTime": "2026-02-22T18:00:00+01:00",
    "venue": "Felsenkeller",
    "category": "concert",
    "price": "12€",
    "url": "https://...",
    "source": "leipzig-im"
  }
]
```

### Compact (one-liner)
```
Sat 22 Feb 18:00 | 🎵 Jinjer – European Duél Tour @ Felsenkeller (12€)
Sat 22 Feb 10:00 | 🛍️ Flohmarkt Plagwitz @ Markthalle (free)
```

## Data Sources (Phase Plan)

### Phase 1 — MVP
**leipzig.de** — City's official event calendar
- URL: `https://www.leipzig.de/kultur-und-freizeit/veranstaltungen`
- Categories: Bühne, Freizeit, Kinder & Jugendliche, Konzert, Lesung, Ausstellungen
- Has date range and category query params
- Server-rendered HTML
- Broadest coverage including family events, civic events, markets

### Phase 2 — Scene Coverage
**leipzig-im.de** — Local scene event calendar
- URL: `https://www.leipzig-im.de/index.php?section=anzeigensort&sort=rubrik`
- Categories: Konzert, Kinder & Familie, Theater, Ausstellungen, Sport, Märkte
- All the indie venues (Werk 2, UT Connewitz, Felsenkeller, Täubchenthal, Conne Island)
- Server-rendered HTML

### Phase 3 — API Source
**Songkick** — Concert aggregator with API
- URL: `https://www.songkick.com/metro-areas/28528-germany-leipzig`
- 748+ events, music-focused
- Has a public API (needs key)
- Good for touring acts that local sites might miss

### Future
- **livegigs.de** — German concert aggregator
- **Bandsintown** — Artist-focused API
- **Individual venue sites** (Täubchenthal, Felsenkeller, etc.)

## Technical Decisions

| Decision | Choice | Why |
|----------|--------|-----|
| Language | Go | Consistent with bahn-cli, zoo-leipzig. Fast, single binary, great for CLI + scraping. |
| CLI framework | Cobra | Same as other projects. Proven. |
| HTTP | net/http + colly or goquery | goquery for HTML parsing, standard lib for HTTP |
| Output | JSON to stdout | Agent-first. Human formats are sugar on top. |
| Diagnostics | stderr | Errors, warnings, progress — never pollute stdout |
| Cache | SQLite | Local cache for scraped data. Avoid hammering sources. TTL per source. |
| Config | `~/.config/leipzig/config.yaml` | Source enable/disable, API keys, cache TTL |

## Caching Strategy
- Each source's results cached in SQLite with configurable TTL
- Default TTL: 1 hour (events don't change that fast)
- `leipzig cache clear` to force refresh
- Cache key: source + date range hash
- Stale cache served if source is unreachable (with warning on stderr)

## Deduplication
Same event may appear on multiple sources. Dedupe by:
1. Normalize event name (lowercase, strip punctuation)
2. Match on name similarity + same date + same venue
3. Keep the version with more detail (longer description, has price, etc.)
4. Mark with all sources it appeared in

## Project Structure

```
leipzig-cli/
├── cmd/
│   ├── root.go           # Cobra root command
│   ├── events.go         # `leipzig events` command
│   ├── sources.go        # `leipzig sources` command
│   └── cache.go          # `leipzig cache` command
├── internal/
│   ├── model/
│   │   └── event.go      # Canonical Event type + Category constants
│   ├── source/
│   │   ├── source.go     # Source interface definition
│   │   ├── leipzigde/    # leipzig.de adapter
│   │   ├── leipzigim/    # leipzig-im.de adapter
│   │   └── songkick/     # songkick adapter
│   ├── engine/
│   │   ├── engine.go     # Orchestrates sources, merge, dedupe
│   │   └── filter.go     # Filtering logic
│   ├── cache/
│   │   └── cache.go      # SQLite cache layer
│   └── output/
│       ├── json.go       # JSON formatter
│       ├── table.go      # Table formatter
│       └── compact.go    # Compact formatter
├── main.go
├── go.mod
├── go.sum
├── SSD.md
└── README.md
```

## Non-Goals (for now)
- No ticket purchasing / booking
- No user accounts or favorites
- No push notifications (Ori handles that via cron)
- No GUI / web interface
- No venue discovery (focus is events, not "find me a bar")

## Integration with Ori
Once built, Ori can:
- Run `leipzig events --weekend --json` every Friday afternoon
- Filter by weather (outdoor events only if sunny)
- Cross-reference with calendar (skip conflicts)
- Suggest family activities on days Elio isn't at Kita
- Pair concert discoveries with curator for pre-event playlists

## Open Questions
1. **Rate limiting:** How aggressive can we scrape leipzig.de / leipzig-im.de? Should we add configurable delays?
2. **Geo data:** Worth adding lat/lng to venues for distance-based filtering? (Nice for "stuff near us")
3. **Image URLs:** Some events have poster images — include in model or skip?

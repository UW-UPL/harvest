# Technical Details

Aggregates RSS feeds from UPL member blogs into a unified feed. Handles a bunch of RSS formats because everyone implements them differently:

### Core Fields

The system handles both RSS and Atom feeds with different structures:

**RSS Items:**
```go
type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    PubDate     string `xml:"pubDate"`
    Author      string `xml:"author"`
    Creator     string `xml:"dc:creator"`
    Description string `xml:"description"`
    Content     string `xml:"content:encoded"` // full content
}
```

**Atom Entries:**
```go
type AtomEntry struct {
    Title     string     `xml:"title"`
    Links     []AtomLink `xml:"link"`         // rel="alternate" for the actual link
    Published string     `xml:"published"`
    Updated   string     `xml:"updated"`      // fallback date
    Author    AtomAuthor `xml:"author"`       // nested author.name
    Content   string     `xml:"content"`
    Summary   string     `xml:"summary"`
}
```

**Unified Output:**
```go
type BlogPost struct {
    Title   string
    Link    string
    Date    time.Time
    Author  string    // falls back to channel/feed title if missing
    Summary string    // cleaned and truncated to 200 chars
}
```

### Date Hell
RSS feeds use whatever date format they feel like. We handle:
```go
var dateFormats = []string{
    time.RFC1123Z,                        // Mon, 02 Jan 2006 15:04:05 -0700
    time.RFC1123,                         // Mon, 02 Jan 2006 15:04:05 MST
    time.RFC3339,                         // 2006-01-02T15:04:05Z07:00 (atom's favorite)
    time.RFC3339Nano,                     // with nanoseconds
    "2006-01-02T15:04:05Z",               // simplified ISO
    "2006-01-02 15:04:05 -0700",          // space instead of T
    "02 Jan 2006 15:04 -0700",            // why do people use this
    "Mon, 02 Jan 2006 15:04:05 GMT",      // GMT variant
    "02 Jan 2006 15:04 +0000",            // another variant
    "2006-01-02",                         // at least it's simple
    "January 2, 2006",                    // human readable
}
```

### Content Cleanup
- Strips HTML tags with regex
- Unescapes HTML entities (`&amp;` → `&`)
- Normalizes whitespace
- Truncates to 200 characters max
- Strips markdown characters to avoid formatting issues
- Falls back to "Visit post for details." if no content available

## How It Works

1. Reads feeds from `whitelist.toml`
2. Downloads them all concurrently using goroutines (because waiting sucks)
3. Parses XML, handles both RSS and Atom formats
4. Extracts and normalizes:
   - **Title**: Direct mapping
   - **Link**: Direct for RSS, finds `rel="alternate"` for Atom
   - **Author**: Tries multiple fields, falls back to channel/feed title
   - **Date**: Tries multiple date fields and formats
   - **Summary**: Cleans content/description, truncates to 200 chars
5. Sorts posts by date (newest first)
6. Generates two output formats:
   - JSON feed at `output/blog_posts.json`
   - RSS/XML feed at `output/feed.xml`

The code's modular so we can add new formats when someone inevitably implements RSS wrong again. This has been so fun to troubleshoot.

# Technical Details

## RSS Support

Handles a bunch of RSS formats because everyone implements them differently:

### Core Fields
```go
type Item struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    PubDate     string `xml:"pubDate"`      // regular rss
    Date        string `xml:"date"`         // some use this
    Published   string `xml:"published"`    // atom folks
    Updated     string `xml:"updated"`      // fallback
}
```

### Date Hell
RSS feeds use whatever date format they feel like. We handle:
```go
var dateFormats = []string{
    time.RFC1123Z,              // most RSS
    time.RFC3339,               // atom's favorite
    "02 Jan 2006 15:04 -0700",  // why do people use this
    "2006-01-02",               // at least it's simple
}
```

### Content Cleanup
- Strips HTML (nobody needs that in a feed)
- Fixes entities (`&amp;` → `&`)
- Handles missing descriptions (minimalist blogs)

## How It Works

1. Reads feeds from `whitelist.toml`
2. Downloads them all at once (because waiting sucks)
3. Parses XML, prays it's valid
4. Cleans up the mess
5. Dumps a nice markdown file in `output/`

The code's modular so we can add new formats when someone inevitably implements RSS wrong again. This has been so fun to troubleshoot.

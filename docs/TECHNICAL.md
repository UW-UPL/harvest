Technical Details
=================

harvest aggregates RSS and Atom feeds from UPL member blogs into a
single unified output. In practice, that means handling a fair bit of
format variation, since not every publisher implements RSS the same
way.

Each RSS item is parsed into the following struct:

    type RSSItem struct {
        Title       string `xml:"title"`
        Link        string `xml:"link"`
        PubDate     string `xml:"pubDate"`
        Author      string `xml:"author"`
        Creator     string `xml:"dc:creator"`
        Description string `xml:"description"`
        Content     string `xml:"content:encoded"`
    }

Atom entries have a slightly different shape: links are an array whose
`rel="alternate"` entry points at the post, and the author sits inside
a nested element.

    type AtomEntry struct {
        Title     string     `xml:"title"`
        Links     []AtomLink `xml:"link"`
        Published string     `xml:"published"`
        Updated   string     `xml:"updated"`
        Author    AtomAuthor `xml:"author"`
        Content   string     `xml:"content"`
        Summary   string     `xml:"summary"`
    }

Both are reduced to a common `BlogPost` record carrying a title, link,
timestamp, author (falling back to the feed title when absent), and a
summary that has been cleaned and truncated to 200 characters.

Dates are the hardest part of the job. Feeds arrive in a stew of
RFC1123, RFC1123Z, RFC3339, and a long tail of ad-hoc variants, so the
parser tries each of the following in order and stops on the first
that succeeds:

    var dateFormats = []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC3339,
        time.RFC3339Nano,
        "2006-01-02T15:04:05Z",
        "2006-01-02 15:04:05 -0700",
        "02 Jan 2006 15:04 -0700",
        "Mon, 02 Jan 2006 15:04:05 GMT",
        "02 Jan 2006 15:04 +0000",
        "2006-01-02",
        "January 2, 2006",
    }

Content cleanup strips HTML tags, unescapes entities, normalizes
whitespace, drops stray markdown, and truncates the result to 200
characters. When a feed provides nothing usable, the summary falls
back to "Visit post for details."

At the top level, harvest reads `whitelist.toml`, fetches every listed
feed concurrently, parses RSS and Atom, normalizes each entry, sorts
the combined list by date (newest first), and writes two outputs:
`output/blog_posts.json` and `output/feed.xml`. The parser is modular
so new formats can be slotted in when the next idiosyncratic feed
turns up.

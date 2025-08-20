package feed

import (
	"encoding/xml"
	"time"
)

type Config struct {
	Feeds []string `toml:"feeds"`
}

// RSS Item structure
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
	Creator     string `xml:"dc:creator"`
	Description string `xml:"description"`
	Content     string `xml:"content:encoded"`
}

// Atom structures
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}

type AtomEntry struct {
	Title     string     `xml:"title"`
	Links     []AtomLink `xml:"link"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
	Author    AtomAuthor `xml:"author"`
	Content   string     `xml:"content"`
	Summary   string     `xml:"summary"`
}

// Unified Item structure for internal use
type Item struct {
	Title       string
	Link        string
	PubDate     string
	Date        string
	Published   string
	Updated     string
	Author      string
	Creator     string
	Description string
	Content     string
	Encoded     string
	Summary     string
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type Feed struct {
	XMLName xml.Name    `xml:""`
	Channel Channel     `xml:"channel"`
	Entries []AtomEntry `xml:"entry"`
	Title   string      `xml:"title"`
}

type BlogPost struct {
	Title   string
	Link    string
	Date    time.Time
	Author  string
	Summary string
}
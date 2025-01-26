package feed

import "time"

type Config struct {
	Feeds []string `toml:"feeds"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Date        string `xml:"date"`
	Published   string `xml:"published"`
	Updated     string `xml:"updated"`
	Author      string `xml:"author"`
	Creator     string `xml:"creator"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
	Encoded     string `xml:"encoded"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
	Entries     []Item `xml:"entry"`
}

type Feed struct {
	Channel Channel `xml:"channel"`
	Entries []Item  `xml:"entry"`
}

type BlogPost struct {
	Title   string
	Link    string
	Date    time.Time
	Author  string
	Summary string
}

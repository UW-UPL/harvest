package rss

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/UW-UPL/harvest/src/feed"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language,omitempty"`
	PubDate     string    `xml:"pubDate,omitempty"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Author      string `xml:"author,omitempty"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func GenerateRSSFeed(posts []feed.BlogPost, outputPath string) error {
	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       "UPL Member Blogs",
			Link:        "https://upl.cs.wisc.edu",
			Description: "Aggregated blog posts from UPL members",
			Language:    "en-us",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       make([]RSSItem, 0, len(posts)),
		},
	}

	for _, post := range posts {
		item := RSSItem{
			Title:       post.Title,
			Link:        post.Link,
			Description: post.Summary,
			Author:      post.Author,
			PubDate:     post.Date.Format(time.RFC1123Z),
			GUID:        post.Link,
		}
		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	if len(rss.Channel.Items) > 0 {
		rss.Channel.PubDate = rss.Channel.Items[0].PubDate
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling RSS: %w", err)
	}

	xmlHeader := []byte(xml.Header)
	fullOutput := append(xmlHeader, output...)

	if err := os.WriteFile(outputPath, fullOutput, 0644); err != nil {
		return fmt.Errorf("writing RSS file: %w", err)
	}

	return nil
}
package feed

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

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

// we're rendering to markdown so to preserve formatting we need to strip out any markdown characters
func stripMarkdown(input string) string {
	invalidChars := []string{"*", "_", "#", "`", ">", "<", "[", "]", "(", ")", "!", "~", "|", "{", "}", "+"}
	for _, char := range invalidChars {
		input = strings.ReplaceAll(input, char, "")
	}

	return input
}

func cleanHTML(input string, maxLength int) string {
	// first, remove HTML tags
	tagRegex := regexp.MustCompile("<[^>]*>")
	cleaned := tagRegex.ReplaceAllString(input, "")

	// & convert HTML entities
	cleaned = html.UnescapeString(cleaned)

	// & normalize whitespace
	wsRegex := regexp.MustCompile(`\s+`)
	cleaned = wsRegex.ReplaceAllString(cleaned, " ")

	// & truncate to maxLength
	if len(cleaned) > maxLength {
		cleaned = cleaned[:maxLength]

		if cleaned[len(cleaned)-1] == ' ' || cleaned[len(cleaned)-1] == '.' {
			cleaned = cleaned[:len(cleaned)-1]
		}

		cleaned += "..."
	}

	return strings.TrimSpace(cleaned)
}

func parseDate(dateStrs ...string) time.Time {
	for _, dateStr := range dateStrs {
		if dateStr == "" {
			continue
		}

		for _, format := range dateFormats {
			if t, err := time.Parse(format, dateStr); err == nil {
				return t
			}
		}
	}

	log.Printf("warn: Could not parse any date")
	return time.Now()
}

func getDescription(candidates ...string) string {
	for _, candidate := range candidates {
		if candidate != "" {
			return stripMarkdown(cleanHTML(candidate, 200))
		}
	}
	return "Visit post for details."
}

func FetchFeed(url string) ([]BlogPost, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching feed %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading res from %s: %w", url, err)
	}

	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parsing feed %s: %w", url, err)
	}

	var posts []BlogPost

	// handle RSS feeds
	if len(feed.Channel.Items) > 0 {
		for _, item := range feed.Channel.Items {
			author := strings.TrimSpace(item.Author)
			if author == "" {
				author = strings.TrimSpace(item.Creator)
			}
			if author == "" {
				author = feed.Channel.Title
			}

			post := BlogPost{
				Title:   strings.TrimSpace(item.Title),
				Link:    item.Link,
				Date:    parseDate(item.PubDate),
				Author:  author,
				Summary: getDescription(item.Description, item.Content),
			}
			posts = append(posts, post)
		}
	}

	// handle Atom feeds
	if len(feed.Entries) > 0 {
		channelTitle := feed.Title

		for _, entry := range feed.Entries {
			// find the alternate link
			link := ""
			for _, l := range entry.Links {
				if l.Rel == "alternate" || l.Rel == "" {
					link = l.Href
					break
				}
			}

			author := strings.TrimSpace(entry.Author.Name)
			if author == "" {
				author = channelTitle
			}

			post := BlogPost{
				Title:   strings.TrimSpace(entry.Title),
				Link:    link,
				Date:    parseDate(entry.Published, entry.Updated),
				Author:  author,
				Summary: getDescription(entry.Summary, entry.Content),
			}
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func FetchAllFeeds(feeds []string) []BlogPost {
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		posts []BlogPost
	)

	for _, feedURL := range feeds {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			feedPosts, err := FetchFeed(url)
			if err != nil {
				log.Printf("err fetching %s: %v", url, err)
				return
			}

			mu.Lock()
			posts = append(posts, feedPosts...)
			mu.Unlock()
		}(feedURL)
	}

	wg.Wait()

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts
}

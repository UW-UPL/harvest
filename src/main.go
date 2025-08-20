package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/UW-UPL/harvest/src/feed"
	"github.com/UW-UPL/harvest/src/json"
	"github.com/UW-UPL/harvest/src/rss"
)

func readConfig(path string) (feed.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return feed.Config{}, fmt.Errorf("reading config: %w", err)
	}

	log.Printf("read config data: %s", string(data))

	var config feed.Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return feed.Config{}, fmt.Errorf("parsing config: %w", err)
	}

	log.Printf("parsed feeds: %v", config.Feeds)
	return config, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	config, err := readConfig("whitelist.toml")
	if err != nil {
		log.Fatalf("err reading config: %v", err)
	}

	posts := feed.FetchAllFeeds(config.Feeds)
	log.Printf("Fetched %d posts", len(posts))

	if err := json.Generate(posts, "output/blog_posts.json"); err != nil {
		log.Fatalf("err generating json: %v", err)
	}

	if err := rss.GenerateRSSFeed(posts, "output/feed.xml"); err != nil {
		log.Fatalf("err generating RSS feed: %v", err)
	}
	log.Printf("Generated RSS feed at output/feed.xml")
}

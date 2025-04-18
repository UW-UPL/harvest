package markdown

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/UW-UPL/harvest/src/feed"
)

// struct
type postOut struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Author      string `json:"author"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

func Generate(posts []feed.BlogPost, outputPath string) error {
	// make sure the directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	// build the slice we want to dump
	out := struct {
		Posts []postOut `json:"posts"`
	}{}

	for _, p := range posts {
		out.Posts = append(out.Posts, postOut{
			Title:       p.Title,
			Link:        p.Link,
			Author:      p.Author,
			Date:        p.Date.Format(time.DateOnly), // 2006‑01‑02
			Description: p.Summary,
		})
	}

	// pretty‑print JSON straight to the file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}

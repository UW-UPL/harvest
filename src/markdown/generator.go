package markdown

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/UW-UPL/harvest/src/feed"
)

func Generate(posts []feed.BlogPost, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "# UPL Blog Posts\n")

	for _, post := range posts {
		fmt.Fprintf(file, "## [%s](%s)\n", post.Title, post.Link)
		fmt.Fprintf(file, "*By %s on %s*\n\n", post.Author, post.Date.Format("2006-01-02"))
		fmt.Fprintf(file, "%s\n\n---\n\n", post.Summary)
	}

	return nil
}

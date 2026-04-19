Contributing
============

To add your blog, fork the repository, append your RSS feed URL to
`whitelist.toml`, and open a pull request. Feeds should belong to a
UPL member's personal blog, expose proper RSS or Atom (not a sitemap),
and include at minimum a title, link, and publication date on each
entry. Posts are expected to be tech-focused.

Before opening the PR, it is worth confirming the project still builds
and runs against the updated whitelist:

    go mod tidy
    go build -o harvest ./src
    ./harvest

If you run into anything odd while adding a feed — unusual date
formats, missing authors, and so on — [TECHNICAL.md](TECHNICAL.md)
describes what the parser accepts and where the normalization
happens. The project
layout is kept intentionally simple; please keep it that way.

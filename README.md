harvest
=======

harvest is a small Go program that aggregates posts from the RSS feeds
of UPL member blogs into a single `output/blog_posts.json` and
`output/feed.xml`. Feeds are defined in [`whitelist.toml`](whitelist.toml),
and a GitHub Actions workflow rebuilds the outputs every 12 hours, 
but only when one or more whitelisted blogs have changed.

To add your blog, open a pull request against `whitelist.toml` (see
[`docs/CONTRIBUTING.md`](docs/CONTRIBUTING.md). To run it locally, `go build -o harvest ./src`
followed by `./harvest` will fetch the listed feeds and write to
`output/`.

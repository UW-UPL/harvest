# Contributing

1. Fork this repository
2. Add your RSS feed URL to `whitelist.toml`
3. Create a PR

Requirements:
- Must be a UPL member's personal blog RSS feed
- No sitemaps
- Feed must include title, link, and publication date
- Posts should be tech-focused

```bash
go mod tidy
go build -o harvest ./src
./harvest
```

The file structure is self-explanatory. Please keep it that way! :]

# Contributing to Harvest

## Adding Your Blog

1. Fork this repository
2. Add your RSS feed URL to `whitelist.toml`
3. Create a PR

Requirements:
- Must be a personal blog RSS feed
- No sitemaps
- Feed must include title, link, and publication date
- Posts should be tech-focused

## Development

### Setup
```bash
go mod tidy
```

### Build
```bash
go build -o harvest ./src
```

### Run
```bash
./harvest
```

### Structure
- `src/`: Go source code
- `output/`: Generated blog post markdown
- `whitelist.toml`: RSS feed list

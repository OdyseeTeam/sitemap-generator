# Odysee Sitemap Generator
[![Build Status](https://travis-ci.com/OdyseeTeam/sitemap-generator.svg?branch=master)](https://travis-ci.com/OdyseeTeam/sitemap-generator)

This tool builds sitemaps for https://odysee.com

# Requirements
- A Chainquery mysql connection (credentials or own instance)
- Quite some RAM (works with 32GB for sure)
- A web server if you want to serve the sitemaps

# Setup
- Get the latest release from the releases page on GitHub
- Create and fill `config.json` using [this example](config.json.example)

# Instructions

```
builds and uploads a sitemap for odysee.com

Usage:
  sitemap-generator [flags]

Flags:
      --embed-endpoint string     endpoint for embeds (default "https://odysee.com/$/embed/")
  -h, --help                      help for sitemap-generator
      --player-endpoint string    endpoint of the player (default "https://cdn.lbryplayer.xyz/")
      --sitemap-endpoint string   endpoint for embeds (default "https://sitemaps.odysee.com")
      --website string            endpoint for embeds (default "https://odysee.com")
```

## Running from Source

Clone the repository and run `make`

## License

This project is MIT licensed. For the full license, see [LICENSE](LICENSE).

## Contributing

Contributions to this project are welcome, encouraged, and compensated.

## Security

We take security seriously. Please contact [security@odysee.com](mailto:security@odysee.com) regarding any security issues.

## Contact

The primary contact for this project is [Niko Storni](https://github.com/nikooo777) (niko@odysee.com).

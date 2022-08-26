# Odysee Sitemap Generator
[![Build Status](https://app.travis-ci.com/OdyseeTeam/sitemap-generator.svg?token=6GPvhdnmp2bbM55sBtou&branch=master)](https://app.travis-ci.com/OdyseeTeam/sitemap-generator)
[![Latest release](https://badgen.net/github/release/OdyseeTeam/sitemap-generator)](https://github.com/OdyseeTeam/sitemap-generator/releases)

This tool builds sitemaps for https://odysee.com

# Requirements
- A Chainquery mysql connection (credentials or own instance)
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
      --player-endpoint string    endpoint of the player (default "https://player.odycdn.com/")
      --sitemap-endpoint string   endpoint for embeds (default "https://sitemaps.odysee.com")
      --type string               type of sitemap to generate (global, monthly, weekly, daily, hourly (default "global")
      --website string            endpoint for embeds (default "https://odysee.com")
```

## Running from Source

Clone the repository and run `make`

## Contributing

Contributions to this project are welcome, encouraged, and compensated.

## Security

We take security seriously. Please contact [security@odysee.com](mailto:security@odysee.com) regarding any security issues.

## Contact

The primary contact for this project is [Niko Storni](https://github.com/nikooo777) (niko@odysee.com).

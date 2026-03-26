# Wikifreedia Relay

A Nostr relay implementation for wiki events, built with [Khatru](https://github.com/fiatjaf/khatru) and Go.

## Features

- Accepts only specific wiki-related event kinds (30818, 30819, 818, 819)
- SQLite3 backend for event storage
- Fast and efficient Nostr relay implementation
- Full-text search support for wiki events

## Building

```bash
go build -o relay
```

## Running

```bash
./relay
```

The relay will start on `0.0.0.0:3334` by default.

## Configuration

- Event database is stored in `./db` (SQLite3)
- Only event kinds 30818, 30819, 818, 819 are accepted
- Relay name: "Wikifreedia relay"
- Relay pubkey: `fa984bd7dbb282f07e16e7ae87b26a2a7b9b90b7246a44771f0cf5ae58018f52`

## Development

Dependencies:
- `github.com/fiatjaf/khatru` - Nostr relay framework
- `github.com/fiatjaf/eventstore` - Event storage interface
- `github.com/nbd-wtf/go-nostr` - Nostr utilities
- `github.com/mattn/go-sqlite3` - SQLite3 driver

## License

MIT

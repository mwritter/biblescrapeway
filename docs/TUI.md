# External TUI / automation contract

This binary is intended to be driven by another program (for example a TUI) as a subprocess.

## Invoking the binary

Build from the module root:

```bash
go build -o scraper ./cmd/scraper
```

Example:

```bash
./scraper -passage 'John 3:16' -version NIV -json
```

Use `-cache-dir /path/to/dir` if the caller should use an isolated cache (for example per-session or per-user profile in the TUI).

## Standard streams

- **stdout**: On success, either human-readable verse lines (default) or a **single JSON object** when `-json` is set. Do not parse stdout as JSON unless `-json` was passed **and** the process exited with code 0.
- **stderr**: Diagnostics. With `-verbose`, the program logs cache behavior, the scraping URL, and network/cache lifecycle messages via the standard library `log` package (timestamps prefixed). Without `-verbose`, stderr is quiet except for errors emitted by `run()` on failure.

## Exit codes

- **0**: Success. With `-json`, stdout contains one JSON document matching [`scraper.Entry`](../internal/scraper/cache.go) (`fetched_at`, `passages` with `book`, `chapter`, `verse`, `text`, `version`).
- **1**: Failure (flags, cache directory resolution, scrape error, cache write error, JSON encode error). An error line is written to stderr. **Do not** treat stdout as JSON on non-zero exit.

## JSON shape

The success payload is stable for tooling:

```json
{
  "fetched_at": "2026-05-15T12:00:00Z",
  "passages": [
    {
      "book": "John",
      "chapter": 3,
      "verse": 16,
      "text": "...",
      "version": "NIV"
    }
  ]
}
```

## Importing as a library (Go)

The same logic lives in `scraper/internal/scraper`. A Go TUI can call `scraper.Fetch` with `scraper.Options` instead of shelling out, if preferred.

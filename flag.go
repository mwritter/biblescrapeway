package main

import "flag"

type Flags struct {
	Passage  string
	Version  string
	CacheDir string
	Refresh  bool
	Json     bool
	Verbose  bool
}

func parseFlags() (*Flags, error) {
	passage := flag.String("passage", "john 3:16-20, gen 1:1", "Bible Gateway passage search text (trimmed; internal whitespace collapsed; version uppercased for cache key)")
	version := flag.String("version", "NIV", "Bible translation code (uppercased for cache key)")
	cacheDirFlag := flag.String("cache-dir", "", "Directory for JSON cache files (default: $XDG_CACHE_HOME/scraper or ~/.cache/scraper)")
	refresh := flag.Bool("refresh", false, "Ignore cache, fetch from the network, and replace the cache entry on success (updates fetched_at)")
	jsonOut := flag.Bool("json", false, "Print result as JSON to stdout on success (see docs/TUI.md). Diagnostics use stderr.")
	verbose := flag.Bool("verbose", false, "Log cache hit, miss, bypass, read errors, network fetch, and scraping URL (all to stderr)")
	flag.Parse()
	return &Flags{
		Passage:  *passage,
		Version:  *version,
		CacheDir: *cacheDirFlag,
		Refresh:  *refresh,
		Json:     *jsonOut,
		Verbose:  *verbose,
	}, nil
}

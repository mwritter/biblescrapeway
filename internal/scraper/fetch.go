package scraper

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Options configures Fetch.
type Options struct {
	Passage  string
	Version  string
	CacheDir string
	Refresh  bool
	Verbose  bool
}

// Fetch resolves a passage (cache hit or network scrape) and returns an entry.
// CacheDir must be non-empty (use DefaultCacheDir if needed).
// Diagnostics when Verbose is true go to the standard logger (stderr by default).
func Fetch(ctx context.Context, opts Options) (Entry, error) {
	if err := ctx.Err(); err != nil {
		return Entry{}, err
	}
	if opts.CacheDir == "" {
		return Entry{}, fmt.Errorf("scraper: CacheDir is empty")
	}

	pNorm, vNorm := normalizeCacheInputs(opts.Passage, opts.Version)
	var entry Entry
	fromCache := false

	if !opts.Refresh {
		e, ok, err := loadCache(opts.CacheDir, pNorm, vNorm)
		if err != nil {
			if opts.Verbose {
				log.Printf("cache: %v", err)
			}
		} else if ok {
			entry = e
			fromCache = true
			if opts.Verbose {
				log.Printf("cache hit (fetched_at=%s)", entry.FetchedAt.UTC().Format(time.RFC3339))
			}
		} else if opts.Verbose {
			log.Print("cache miss (no entry)")
		}
	} else if opts.Verbose {
		log.Print("cache bypass (-refresh)")
	}

	if !fromCache {
		if err := ctx.Err(); err != nil {
			return Entry{}, err
		}
		sc := scraperImpl{}
		if err := sc.scrapePassage(pNorm, vNorm, opts.Verbose); err != nil {
			return Entry{}, fmt.Errorf("scrape: %w", err)
		}
		entry = Entry{
			FetchedAt: time.Now().UTC(),
			Passages:  sc.Passages,
		}
		if err := saveCache(opts.CacheDir, pNorm, vNorm, entry); err != nil {
			return Entry{}, fmt.Errorf("cache write: %w", err)
		}
		if opts.Verbose {
			log.Printf("fetched from network; wrote cache (fetched_at=%s)", entry.FetchedAt.UTC().Format(time.RFC3339))
		}
	}

	return entry, nil
}

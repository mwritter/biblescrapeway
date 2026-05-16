package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"scraper/internal/scraper"
)

func main() {
	log.SetOutput(os.Stderr)
	if err := run(); err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}

func run() error {
	flags, err := parseFlags()
	if err != nil {
		return fmt.Errorf("flags: %w", err)
	}

	cacheDir := flags.CacheDir
	if cacheDir == "" {
		d, err := scraper.DefaultCacheDir()
		if err != nil {
			return fmt.Errorf("cache directory: %w", err)
		}
		cacheDir = d
	}

	entry, err := scraper.Fetch(context.Background(), scraper.Options{
		Passage:  flags.Passage,
		Version:  flags.Version,
		CacheDir: cacheDir,
		Refresh:  flags.Refresh,
		Verbose:  flags.Verbose,
	})
	if err != nil {
		return err
	}

	if flags.Json {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(entry); err != nil {
			return fmt.Errorf("json: %w", err)
		}
		return nil
	}

	for _, p := range entry.Passages {
		fmt.Println(p.String())
	}
	return nil
}

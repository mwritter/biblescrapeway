package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Entry is the JSON shape for one (passage, version) scrape (stdout with -json, or on-disk cache).
type Entry struct {
	FetchedAt time.Time `json:"fetched_at"`
	Passages  []Passage `json:"passages"`
}

func normalizeCacheInputs(passage, version string) (passageNorm, versionNorm string) {
	p := strings.TrimSpace(passage)
	p = strings.Join(strings.Fields(p), " ")
	v := strings.ToUpper(strings.TrimSpace(version))
	return p, v
}

func cacheKeyMaterial(passage, version string) string {
	p, v := normalizeCacheInputs(passage, version)
	return p + "\x00" + v
}

func cacheFileName(passage, version string) string {
	sum := sha256.Sum256([]byte(cacheKeyMaterial(passage, version)))
	return hex.EncodeToString(sum[:]) + ".json"
}

func cacheFilePath(dir, passage, version string) string {
	return filepath.Join(dir, cacheFileName(passage, version))
}

// loadCache returns a cache entry if the file exists and unmarshals cleanly.
// If the file is missing, it returns ok=false and err=nil.
// If the file is unreadable or invalid JSON, it returns err non-nil.
func loadCache(dir, passage, version string) (e Entry, ok bool, err error) {
	path := cacheFilePath(dir, passage, version)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, false, nil
		}
		return Entry{}, false, err
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, false, fmt.Errorf("cache %s: %w", path, err)
	}
	return e, true, nil
}

func saveCache(dir, passage, version string, e Entry) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := cacheFilePath(dir, passage, version)
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// DefaultCacheDir returns $XDG_CACHE_HOME/scraper when XDG_CACHE_HOME is set,
// otherwise ~/.cache/scraper.
func DefaultCacheDir() (string, error) {
	if d := os.Getenv("XDG_CACHE_HOME"); d != "" {
		return filepath.Join(d, "scraper"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "scraper"), nil
}

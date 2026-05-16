package scraper

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNormalizeCacheInputs(t *testing.T) {
	tests := []struct {
		passage, version     string
		wantPassage, wantVer string
	}{
		{"  john 3:16  ", "niv", "john 3:16", "NIV"},
		{"gen 1:1", "  esv ", "gen 1:1", "ESV"},
		{"a  b\tc", "kjv", "a b c", "KJV"},
	}
	for _, tt := range tests {
		gotP, gotV := normalizeCacheInputs(tt.passage, tt.version)
		if gotP != tt.wantPassage || gotV != tt.wantVer {
			t.Errorf("normalizeCacheInputs(%q, %q) = (%q, %q), want (%q, %q)",
				tt.passage, tt.version, gotP, gotV, tt.wantPassage, tt.wantVer)
		}
	}
}

func TestCacheKeyMaterialStable(t *testing.T) {
	a := cacheKeyMaterial("  john 3:16 ", "niv")
	b := cacheKeyMaterial("john 3:16", "NIV")
	if a != b {
		t.Errorf("cache key material should match after normalization: %q vs %q", a, b)
	}
}

func TestCacheFileNameDeterministic(t *testing.T) {
	n1 := cacheFileName("john 3:16", "NIV")
	n2 := cacheFileName("  john  3:16  ", "niv")
	if n1 != n2 {
		t.Errorf("cache file name mismatch: %s vs %s", n1, n2)
	}
	const wantLen = 64 + len(".json") // hex(SHA256) + ".json"
	if len(n1) != wantLen {
		t.Errorf("unexpected cache file name length: got %d want %d for %q", len(n1), wantLen, n1)
	}
}

func TestEntryJSONRoundTrip(t *testing.T) {
	ts := time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
	orig := Entry{
		FetchedAt: ts,
		Passages: []Passage{
			{Book: "John", Chapter: 3, Verse: 16, Text: "For God so loved", Version: "NIV"},
		},
	}
	data, err := json.Marshal(&orig)
	if err != nil {
		t.Fatal(err)
	}
	var got Entry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !got.FetchedAt.Equal(orig.FetchedAt) {
		t.Errorf("FetchedAt: got %v want %v", got.FetchedAt, orig.FetchedAt)
	}
	if len(got.Passages) != 1 || got.Passages[0].Book != "John" {
		t.Errorf("Passages: %+v", got.Passages)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	passage, version := "John 3:16", "NIV"
	want := Entry{
		FetchedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		Passages: []Passage{
			{Book: "John", Chapter: 3, Verse: 16, Text: "x", Version: "NIV"},
		},
	}
	if err := saveCache(dir, passage, version, want); err != nil {
		t.Fatal(err)
	}
	got, ok, err := loadCache(dir, passage, version)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected cache hit")
	}
	if !got.FetchedAt.Equal(want.FetchedAt) {
		t.Errorf("FetchedAt got %v want %v", got.FetchedAt, want.FetchedAt)
	}
	if len(got.Passages) != len(want.Passages) {
		t.Fatalf("passages len got %d want %d", len(got.Passages), len(want.Passages))
	}
	if got.Passages[0] != want.Passages[0] {
		t.Errorf("passage got %+v want %+v", got.Passages[0], want.Passages[0])
	}
}

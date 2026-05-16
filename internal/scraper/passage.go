package scraper

import (
	"fmt"
	"strconv"
	"strings"
)

// Passage is one verse from Bible Gateway passage HTML.
type Passage struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
	Version string `json:"version"`
}

// ParseVerseClassToken parses a Bible Gateway verse class token such as
// "John-3-16", "Gen-1-1", or "1Cor-5-7" into book slug, chapter, and verse.
func ParseVerseClassToken(tok string) (book string, chapter, verse int, ok bool) {
	parts := strings.Split(tok, "-")
	if len(parts) < 3 {
		return "", 0, 0, false
	}
	vStr := parts[len(parts)-1]
	cStr := parts[len(parts)-2]
	verse, err := strconv.Atoi(vStr)
	if err != nil {
		return "", 0, 0, false
	}
	chapter, err = strconv.Atoi(cStr)
	if err != nil {
		return "", 0, 0, false
	}
	book = strings.Join(parts[:len(parts)-2], "-")
	if book == "" {
		return "", 0, 0, false
	}
	return book, chapter, verse, true
}

func (p *Passage) String() string {
	return fmt.Sprintf("%s %d:%d — %s (%s)", p.Book, p.Chapter, p.Verse, p.Text, p.Version)
}

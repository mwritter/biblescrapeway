package scraper

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const defaultBibleGatewayBase = "https://www.biblegateway.com/passage/"

func getPlainText(sel *goquery.Selection) string {
	cl := sel.Clone()
	cl.Find("span.chapternum").Remove()
	cl.Find("sup.crossreference").Remove()
	cl.Find("sup.versenum").Remove()
	cl.Find("sup.footnote").Remove()
	t := strings.TrimSpace(cl.Text())
	return strings.Join(strings.Fields(t), " ")
}

// Verse markers on span.text are a second class like John-3-16 or Phil-4-6 (not
// the full book name from the search box). We match that token and skip h1–h4
// where some translations duplicate the same id on a section heading.
var bgVerseClass = regexp.MustCompile(`^[0-9]*[A-Za-z][A-Za-z0-9]*-\d+-\d+$`)

func verseTokenFromSpan(s *goquery.Selection) (tok string, ok bool) {
	for t := range strings.FieldsSeq(s.AttrOr("class", "")) {
		if t != "text" && bgVerseClass.MatchString(t) {
			return t, true
		}
	}
	return "", false
}

// passageFromVerseSpan treats one Bible Gateway verse node (typically
// `.passage-content p > span` with classes `text` and e.g. `John-3-16`) as a single Passage.
func passageFromVerseSpan(s *goquery.Selection) (Passage, bool) {
	if s.Parents().Filter("h1, h2, h3, h4").Length() > 0 {
		return Passage{}, false
	}
	tok, ok := verseTokenFromSpan(s)
	if !ok {
		return Passage{}, false
	}
	book, ch, vs, ok := ParseVerseClassToken(tok)
	if !ok {
		return Passage{}, false
	}
	text := getPlainText(s)
	if text == "" {
		return Passage{}, false
	}
	return Passage{Book: book, Chapter: ch, Verse: vs, Text: text}, true
}

func createPageURL(passage string, version string) (string, error) {
	parsedURL, err := url.Parse(defaultBibleGatewayBase)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %v", err)
	}
	query := parsedURL.Query()
	query.Set("search", passage)
	query.Set("version", version)
	parsedURL.RawQuery = query.Encode()
	pageURL := parsedURL.String()

	return pageURL, nil
}

type scraperImpl struct {
	Passages []Passage
	c        *colly.Collector
}

func (s *scraperImpl) scrapePassage(passage string, version string, verbose bool) error {
	s.Passages = nil

	pageURL, err := createPageURL(passage, version)
	if err != nil {
		return fmt.Errorf("create page URL: %w", err)
	}

	if verbose {
		log.Printf("scraping URL: %s", pageURL)
	}

	s.c = colly.NewCollector()

	s.c.OnHTML(".passage-content p > span", func(e *colly.HTMLElement) {
		if p, ok := passageFromVerseSpan(e.DOM); ok {
			p.Version = version
			s.Passages = append(s.Passages, p)
		}
	})

	if err := s.c.Visit(pageURL); err != nil {
		return fmt.Errorf("visit URL: %w", err)
	}
	return nil
}

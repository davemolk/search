package search

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/davemolk/fuzzyHelpers"
)

// Search takes in a URL, makes a GET request, and parses the response
// body, printing the results to s.output.
func (s *searcher) Search(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request for %s: %v", url, err)
	}

	h, err := fuzzyHelpers.NewHeaders(
		fuzzyHelpers.WithURL(url),
	)
	if err != nil {
		return fmt.Errorf("error with header creation for %s: %v", url, err)
	}
	headers := h.Headers()
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request for %s: %v", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP response: %d for %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot parse response body: %w", err)
	}

	var parse *query

	switch {
	case strings.HasPrefix(url, "https://search.b"):
		parse = s.brave
	case strings.HasPrefix(url, "https://html.d"):
		parse = s.duck
	case strings.HasPrefix(url, "https://www.mo"):
		parse = s.mojeek
	case strings.HasPrefix(url, "https://lite.qwant"):
		parse = s.qwant
	case strings.HasPrefix(url, "https://bing"):
		parse = s.bing
	case strings.HasPrefix(url, "https://search.y"):
		parse = s.yahoo
	// return error if we didn't hit one of these!
	default:
		return fmt.Errorf("mismatched url, check if one of the search engines has changed")
	}

	doc.Find(parse.itemSelector).Each(func(_ int, g *goquery.Selection) {
		var link string
		if parse.name != "qwant" {
			link, _ = g.Find(parse.linkSelector).Attr("href")
		} else {
			link = g.Find(parse.linkSelector).Text()
		}
		blurb := g.Find(parse.blurbSelector).Text()
		cleanedLink := s.cleanLinks(link)
		cleanedBlurb := s.cleanBlurb(blurb)
		s.print(cleanedBlurb, cleanedLink)
	})
	return nil
}

// cleanBlurb does a bit of tidying up of each input blurb string.
func (s *searcher) cleanBlurb(str string) string {
	cleanB := s.noBlank.ReplaceAllString(str, " ")
	cleanB = strings.TrimSpace(cleanB)
	cleanB = strings.ReplaceAll(cleanB, "\n", "")
	return cleanB
}

// cleanLinks does a bit of tidying up of each input URL string.
func (s *searcher) cleanLinks(str string) string {
	u, err := url.QueryUnescape(str)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("unable to clean %s", str))
		return ""
	}
	if strings.HasPrefix(u, "//duck") {
		// ddg links will sometimes take the following format:
		// //duckduckgo.com/l/?uddg=actualURLHere/&rut=otherStuff
		removePrefix := strings.Split(u, "=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "&rut")
		u = removeSuffix[0]
	}
	if strings.HasPrefix(u, "https://r.search.yahoo.com/") {
		removePrefix := strings.Split(u, "/RU=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "/RK=")
		u = removeSuffix[0]
	}
	return u
}

// print truncates any blurb with a length longer
// than s.length and prints to s.output.
func (s *searcher) print(blurb, link string) {
	if len(blurb) > s.length {
		blurb = blurb[:s.length]
	}
	if s.urls && len(blurb) > 0 {
		fmt.Fprintln(s.output, link)
	}
	fmt.Fprintln(s.output, blurb)
	fmt.Fprintln(s.output)
}

package search

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// seed random number generator to get random user agents
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Search takes in a URL, makes a GET request, and parses the response
// body, printing the results to s.output.
func (s *searcher) Search(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request for %s: %v", url, err)
	}

	req = s.headers(req)

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
	case strings.HasPrefix(url, "https://www.ask"):
		parse = s.ask
	case strings.HasPrefix(url, "https://bing"):
		parse = s.bing
	case strings.HasPrefix(url, "https://search.b"):
		parse = s.brave
	case strings.HasPrefix(url, "https://html.d"):
		parse = s.duck
	case strings.HasPrefix(url, "https://search.y"):
		parse = s.yahoo
	// return error if we didn't hit one of these!
	default:
		return fmt.Errorf("mismatched url, check if one of the search engines has changed")
	}

	doc.Find(parse.itemSelector).Each(func(_ int, g *goquery.Selection) {
		link, _ := g.Find(parse.linkSelector).Attr("href")
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
// bing will sometimes encode the links and I haven't bothered to work
// out how to decode them. Maybe one day...
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
	fmt.Println()
}


func (s *searcher) headers(r *http.Request) *http.Request {
	if rand.Intn(2) == 1 {
		return s.ff(r)
	}
	return s.chrome(r)
}

func (s *searcher) ff(r *http.Request) *http.Request {
	uAgent := s.ffUA()
	r.Header.Set("Host", r.URL.Host)
	r.Header.Set("User-Agent", uAgent)
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("DNT", "1")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Sec-GCP", "1")
	return r
}

func (s *searcher) chrome(r *http.Request) *http.Request {
	uAgent := s.chromeUA()
	r.Header.Set("Host", r.URL.Host)
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Cache-Control", "max-age=0")
	r.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`)
	r.Header.Set("sec-ch-ua-mobile", "?0")
	switch s.osys {
	case "m":
		r.Header.Set("sec-ch-ua-platform", "Macintosh")
	default:
		r.Header.Set("sec-ch-ua-platform", "Windows")
	}
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("User-Agent", uAgent)
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	return r
}

func (s *searcher) ffUA() string {
	var userAgents []string
	switch s.osys {
	case "m":
		userAgents = []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:107.0) Gecko/20100101 Firefox/107.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:104.0) Gecko/20100101 Firefox/104.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:101.0) Gecko/20100101 Firefox/101.0",
		}
	case "l":
		userAgents = []string{
			"Mozilla/5.0 (X11; Linux x86_64; rv:93.0) Gecko/20100101 Firefox/93.0 ",
			"Mozilla/5.0 (X11; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0 ",
			"Mozilla/5.0 (X11; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0 ",
			"Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0 ",
			"Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0 ",
		}
	default:
		userAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
		}
	}
	random := rand.Intn(len(userAgents))
	return userAgents[random]
}

func (s *searcher) chromeUA() string {
	var userAgents []string
	switch s.osys {
	case "m":
		userAgents = []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4692.56 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		}
	case "l":
		userAgents = []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4692.56 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		}
	default:
		userAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		}
	}
	random := rand.Intn(len(userAgents))
	return userAgents[random]
}

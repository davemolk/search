package search

import "fmt"

type query struct {
	base string
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
}

func (s *searcher) createQueries() {
	s.ask = &query{
		base:   "https://www.ask.com/web?q=",
		blurbSelector: "div.PartialSearchResults-item p",
		itemSelector:  "div.PartialSearchResults-item",
		linkSelector:  "a.PartialSearchResults-item-title-link",
		name:          "ask",
	}

	s.bing = &query{
		base:   "https://bing.com/search?q=",
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
		name:          "bing",
	}

	s.brave = &query{
		base:   "https://search.brave.com/search?q=",
		blurbSelector: "div.snippet-content p.snippet-description",
		itemSelector:  "div.fdb",
		linkSelector:  "div.fdb > a.result-header",
		name:          "brave",
	}

	s.duck= &query{
		base:   "https://html.duckduckgo.com/html?q=",
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}

	s.yahoo = &query{
		base:   "https://search.yahoo.com/search?p=",
		blurbSelector: "div.compText",
		itemSelector:  "div.algo",
		linkSelector:  "h3 > a",
		name:          "yahoo",
	}
}

func (s *searcher) formatURL() <-chan string {
	// 5 search engines
	out := make(chan string, len(s.terms) * 5)
	switch {
	case s.exact:
		go func() {
				defer close(out)
				for _, term := range s.terms {
					out <- fmt.Sprintf("%s\"%s+%s\"", s.ask.base, s.search, term)
					out <- fmt.Sprintf("%s\"%s+%s\"", s.bing.base, s.search, term)
					out <- fmt.Sprintf("%s\"%s+%s\"", s.brave.base, s.search, term)
					out <- fmt.Sprintf("%s\"%s+%s\"", s.duck.base, s.search, term)
					out <- fmt.Sprintf("%s\"%s+%s\"", s.yahoo.base, s.search, term)
				}
			}()
	default:
		go func() {
			defer close(out)
			for _, term := range s.terms {
				out <- fmt.Sprintf("%s%s+%s", s.ask.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.bing.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.brave.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.duck.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.yahoo.base, s.search, term)
			}
		}()
	}
	
	return out
}
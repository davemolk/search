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
		base:   fmt.Sprintf("%s%s", "https://www.ask.com/web?q=", s.search),
		blurbSelector: "div.PartialSearchResults-item p",
		itemSelector:  "div.PartialSearchResults-item",
		linkSelector:  "a.PartialSearchResults-item-title-link",
		name:          "ask",
	}

	s.bing = &query{
		base:   fmt.Sprintf("%s%s", "https://bing.com/search?q=", s.search),
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
		name:          "bing",
	}

	s.brave = &query{
		base:   fmt.Sprintf("%s%s", "https://search.brave.com/search?q=", s.search),
		blurbSelector: "div.snippet-content p.snippet-description",
		itemSelector:  "div.fdb",
		linkSelector:  "div.fdb > a.result-header",
		name:          "brave",
	}

	s.duck= &query{
		base:   fmt.Sprintf("%s%s", "https://html.duckduckgo.com/html?q=", s.search),
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}

	s.yahoo = &query{
		base:   fmt.Sprintf("%s%s", "https://search.yahoo.com/search?p=", s.search),
		blurbSelector: "div.compText",
		itemSelector:  "div.algo",
		linkSelector:  "h3 > a",
		name:          "yahoo",
	}
}

func (s *searcher) formatURL() <-chan string {
	// 5 search engines
	out := make(chan string, len(s.terms) * 5)
	go func() {
		defer close(out)
		for _, term := range s.terms {
			out <- fmt.Sprintf("%s+%s\n", s.ask.base, term)
			out <- fmt.Sprintf("%s+%s\n", s.bing.base, term)
			out <- fmt.Sprintf("%s+%s\n", s.brave.base, term)
			out <- fmt.Sprintf("%s+%s\n", s.duck.base, term)
			out <- fmt.Sprintf("%s+%s\n", s.yahoo.base, term)
		}
	}()
	return out
}
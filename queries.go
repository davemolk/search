package search

import "fmt"

type query struct {
	base          string
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
}

func (s *searcher) CreateQueries() {
	s.brave = &query{
		base:          "https://search.brave.com/search?q=",
		blurbSelector: "div.snippet-content p.snippet-description",
		itemSelector:  "div.fdb",
		linkSelector:  "div.fdb > a.result-header",
		name:          "brave",
	}
	s.duck = &query{
		base:          "https://html.duckduckgo.com/html?q=",
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}
	s.mojeek = &query{
		base:          "https://www.mojeek.com/search?q=",
		blurbSelector: "li > p.s",
		itemSelector:  "ul.results-standard > li",
		linkSelector:  "li > a.ob",
		name:          "mojeek",
	}
	s.qwant = &query{
		base:          "https://lite.qwant.com/?q=",
		blurbSelector: "article[class='web result'] > p.desc",
		itemSelector:  "article[class='web result']",
		linkSelector:  "article[class='web result'] > span",
		name:          "qwant",
	}
}

func (s *searcher) FormatURL() <-chan string {
	// 4 search engines
	out := make(chan string, len(s.terms)*4)
	switch {
	case s.exact:
		go func() {
			defer close(out)
			for _, term := range s.terms {
				out <- fmt.Sprintf("%s\"%s+%s\"", s.brave.base, s.search, term)
				out <- fmt.Sprintf("%s\"%s+%s\"", s.duck.base, s.search, term)
				out <- fmt.Sprintf("%s\"%s+%s\"", s.mojeek.base, s.search, term)
				out <- fmt.Sprintf("%s\"%s+%s\"", s.qwant.base, s.search, term)
			}
		}()
	default:
		go func() {
			defer close(out)
			for _, term := range s.terms {
				out <- fmt.Sprintf("%s%s+%s", s.brave.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.duck.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.mojeek.base, s.search, term)
				out <- fmt.Sprintf("%s%s+%s", s.qwant.base, s.search, term)
			}
		}()
	}

	return out
}

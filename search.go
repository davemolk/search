package search

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/davemolk/fuzzyHelpers"
)

type searcher struct {
	// query
	exact       bool
	searchExact bool
	multiExact  bool
	multi       bool
	noTerms     bool
	privacy     bool
	search      string
	terms       []string

	// requests
	client      *http.Client
	concurrency int
	debug       bool
	osys        string
	timeout     int

	// output
	length  int
	noBlank *regexp.Regexp
	urls    bool

	// search engines
	bing   *query
	brave  *query
	duck   *query
	mojeek *query
	qwant  *query
	yahoo  *query

	// other
	input  io.Reader
	output io.Writer
}

type option func(*searcher) error

func NewSearcher(opts ...option) (*searcher, error) {
	s := &searcher{
		input:  os.Stdin,
		output: os.Stdout,
	}
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return &searcher{}, err
		}
	}
	return s, nil
}

func WithInput(input io.Reader) option {
	return func(s *searcher) error {
		if input == nil {
			return fmt.Errorf("input is nil")
		}
		s.input = input
		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(s *searcher) error {
		if output == nil {
			return fmt.Errorf("output is nil")
		}
		s.output = output
		return nil
	}
}

var errHelp = errors.New(`usage:
basic query info
-m  include multiple terms within a single query
	search -s foo -m bar baz => https://search.brave.com/search?q=foo+bar+baz, etc.
	default: false
-n  no additional search terms
	search -s foo => https://search.brave.com/search?q=foo, etc.
	default: false
-p  privacy mode (when true, searches brave, duck duck go, mojeek, and qwant,
	otherwise, searches bing, duck duck go, brave, and yahoo)
	default: true
-s  base search term(s)


exact searching
-e  exact searching for entire query
	search -s foo bar -e => https://search.brave.com/search?q="foo+bar", etc.
	default: false
-me exact matching for additional terms
	search -s foo -me bar baz => https://search.brave.com/search?q=foo+"bar+baz", etc.
    default: false
-se exact matching for search term(s)
	search -s "foo bar" -se baz => https://search.brave.com/search?q="foo+bar"+baz, etc.
	default: false


requests
-c  max number of concurrent requests
	default: 10
-os operating system (used for creating browser headers)
	arguments: any, l, m, or w
	default: w
-t  request timeout, in ms
	default: 5000

output
-l  length of result summary
	default: 500
-u  include result urls in output
	default: true

	
help
-d  print the search url to help debug queries
	default: false
-h  help`)

func FromArgs(args []string) option {
	return func(s *searcher) error {
		fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		// query
		multi := fset.Bool("m", false, "multiple terms")
		noTerms := fset.Bool("n", false, "no additional search terms")
		privacy := fset.Bool("p", true, "privacy mode")
		search := fset.String("s", "", "base search term(s)")
		// exact searching
		exact := fset.Bool("e", false, "exact matching")
		multiExact := fset.Bool("me", false, "exact matching for multiple additional terms")
		searchExact := fset.Bool("se", false, "exact matching for base search term(s)")
		//requests
		concurrency := fset.Int("c", 10, "max number of concurrent requests")
		osys := fset.String("os", "w", "l, m, or w")
		to := fset.Int("t", 5000, "timeout in ms")
		// output
		length := fset.Int("l", 500, "length of blurb")
		urls := fset.Bool("u", true, "print urls")
		// help
		debug := fset.Bool("d", false, "print the search url to help debug queries")
		help := fset.Bool("h", false, "")
		fset.SetOutput(s.output)

		err := fset.Parse(args)
		if err != nil {
			return err
		}
		if *help {
			return errHelp
		}

		err = s.validateTerms(*search)
		if err != nil {
			return err
		}
		*search = strings.ReplaceAll(*search, " ", "+")
		err = s.validateOS(*osys)
		if err != nil {
			return err
		}

		s.concurrency = *concurrency
		s.debug = *debug
		s.exact = *exact
		s.length = *length
		s.multi = *multi
		s.multiExact = *multiExact
		s.noTerms = *noTerms
		s.osys = *osys
		s.privacy = *privacy
		s.search = *search
		s.searchExact = *searchExact
		s.timeout = *to
		s.urls = *urls
		s.client = fuzzyHelpers.NewClient(
			fuzzyHelpers.WithConnections(s.concurrency),
		)

		// no additional search terms
		if s.noTerms {
			return nil
		}
		// in case user forgot
		if s.multiExact {
			s.multi = true
		}
		// get terms from args
		args = fset.Args()
		if len(args) > 0 {
			if s.multi {
				m := strings.Join(args, "+")
				args = []string{m}
			}
			s.terms = args
			return nil
		}
		// get terms from s.input if we didn't get from args
		err = s.readTerms()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return nil
	}
}

func (s *searcher) readTerms() error {
	scan := bufio.NewScanner(s.input)
	if s.multi {
		scan.Split(bufio.ScanLines)
	} else {
		scan.Split(bufio.ScanWords)
	}
	for scan.Scan() {
		text := scan.Text()
		if s.multi {
			text = strings.ReplaceAll(scan.Text(), " ", "+")
		}
		s.terms = append(s.terms, text)
	}
	return scan.Err()
}

func RunCLI() {
	s, err := NewSearcher(
		FromArgs(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s.noBlank = regexp.MustCompile(`\s{2,}`)
	s.CreateQueries()
	ch := s.FormatURL()

	tokens := make(chan struct{}, s.concurrency)
	var wg sync.WaitGroup
	for c := range ch {
		wg.Add(1)
		tokens <- struct{}{}
		go func(c string) {
			defer wg.Done()
			defer func() { <-tokens }()
			s.Search(c)
		}(c)
		if s.debug {
			fmt.Println("*****")
			fmt.Println("query:", c)
			fmt.Println("*****")
			fmt.Println()
		}
	}
	wg.Wait()
}

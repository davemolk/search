package search

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type searcher struct {
	// flags
	exact bool
	gophers int
	length int
	multi bool
	osys string
	search string
	terms []string
	timeout int 
	urls bool
	// search engines
	ask *query
	bing *query
	brave *query
	duck *query
	yahoo *query
	// other
	input io.Reader
	noBlank  *regexp.Regexp
	output io.Writer
}

type option func(*searcher) error

func NewSearcher(opts ...option) (*searcher, error) {
	s := &searcher{
		input: os.Stdin,
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

func FromArgs(args []string) option {
	return func(s *searcher) error {
		fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		search := fset.String("s", "", "search term")
		exact := fset.Bool("e", false, "exact matching")
		multi := fset.Bool("m", false, "multiple terms")
		osys := fset.String("os", "w", "m, l, or w")
		to := fset.Int("to", 5000, "timeout in ms")
		urls := fset.Bool("u", false, "print urls")
		length := fset.Int("l", 500, "length of blurb")
		gophers := fset.Int("g", 10, "max number of concurrent requests")
		fset.SetOutput(s.output)

		err := fset.Parse(args)
		if err != nil {
			return err
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

		s.search = *search
		s.exact = *exact
		s.multi = *multi
		s.osys = *osys
		s.timeout = *to
		s.urls = *urls
		s.length = *length
		s.gophers = *gophers

		// get terms from args
		args = fset.Args()
		if len(args) > 1 {
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
	// seed random number generator to get random user agents
	rand.Seed(time.Now().UnixNano())

	s, err := NewSearcher(
		FromArgs(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s.noBlank = regexp.MustCompile(`\s{2,}`)
	s.createQueries()
	ch := s.FormatURL()

	tokens := make(chan struct{}, s.gophers)
	var wg sync.WaitGroup
	for c := range ch {
		wg.Add(1)
		tokens <- struct{}{}
		go func(c string) {
			defer wg.Done()
			s.Search(c)
			<-tokens
		}(c)
	}

	wg.Wait()
}
package search

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type searcher struct {
	input io.Reader
	output io.Writer
	search string
	terms []string
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
		fset.SetOutput(s.output)
		err := fset.Parse(args)
		if err != nil {
			return err
		}
		if *search == "" {
			fmt.Fprintln(os.Stderr, "must provide a search term")
			os.Exit(1)
		}

		s.search = *search

		// get terms
		args = fset.Args()
		if len(args) < 1 {
			return nil
		}
		s.terms = args
		return nil
	}
}

func (s *searcher) readTerms() error {
	scan := bufio.NewScanner(s.input)
	scan.Split(bufio.ScanWords)
	for scan.Scan() {
		s.terms = append(s.terms, scan.Text())
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
	if len(s.terms) < 1 {
		err = s.readTerms()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println(s.search)
	fmt.Println(s.terms)
}
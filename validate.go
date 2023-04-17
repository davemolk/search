package search

import "errors"

var (
	ErrNoSearchTerm = errors.New("must provide search term(s)")
	ErrInvalidOS    = errors.New("os must be l, m, or w")
)

func (s *searcher) validateTerms(str string) error {
	if str == "" {
		return ErrNoSearchTerm
	}
	return nil
}

func (s *searcher) validateOS(str string) error {
	switch str {
	case "l", "m", "w":
		return nil
	default:
		return ErrInvalidOS
	}
}

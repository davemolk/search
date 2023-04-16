package search_test

import (
	"errors"
	"testing"

	"github.com/davemolk/search"
)

func TestNoSearchTerm(t *testing.T) {
	t.Parallel()
	args := []string{}
	_, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if !errors.Is(err, search.ErrNoSearchTerm) {
		t.Fatal("did not fail with ErrNoSearchTerm")
	}
}

func TestInvalidOS(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-os", "bar"}
	_, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if !errors.Is(err, search.ErrInvalidOS) {
		t.Fatal("did not fail with ErrInvalidOS")
	}
}
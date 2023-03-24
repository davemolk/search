package search_test

import (
	"io"
	"testing"

	"github.com/davemolk/search"
)

func TestFromArgsErrorOnBogusFlag(t *testing.T) {
	t.Parallel()
	args := []string{"-bogus"}
	_, err := search.NewSearcher(
		search.WithOutput(io.Discard),
		search.FromArgs(args),
	)
	if err == nil {
		t.Fatal("want error on bogus flag, got nil")
	}
}

/*
func TestFromArgsEmpty(t *testing.T) {
	t.Parallel()
	args := []string{}
	bufInput := bytes.NewBufferString("search terms\nhere")
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
}
*/
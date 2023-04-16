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


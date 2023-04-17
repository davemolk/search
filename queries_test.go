package search_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/davemolk/search"
)

func TestFormatURLSingleTerm(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar")
	args := []string{"-s", "foo"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar",
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermArg(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "bar"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar",
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTerms(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar baz")
	args := []string{"-s", "foo"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar",
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
		"https://www.ask.com/web?q=foo+baz",
		"https://bing.com/search?q=foo+baz",
		"https://search.brave.com/search?q=foo+baz",
		"https://html.duckduckgo.com/html?q=foo+baz",
		"https://search.yahoo.com/search?p=foo+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermArgs(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar",
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
		"https://www.ask.com/web?q=foo+baz",
		"https://bing.com/search?q=foo+baz",
		"https://search.brave.com/search?q=foo+baz",
		"https://html.duckduckgo.com/html?q=foo+baz",
		"https://search.yahoo.com/search?p=foo+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLCombineTerms(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/multi_terms.txt")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-s", "foo", "-m"}
	s, err := search.NewSearcher(
		search.WithInput(f),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar+baz",
		"https://bing.com/search?q=foo+bar+baz",
		"https://search.brave.com/search?q=foo+bar+baz",
		"https://html.duckduckgo.com/html?q=foo+bar+baz",
		"https://search.yahoo.com/search?p=foo+bar+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLCombineTermsArgs(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-m", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar+baz",
		"https://bing.com/search?q=foo+bar+baz",
		"https://search.brave.com/search?q=foo+bar+baz",
		"https://html.duckduckgo.com/html?q=foo+bar+baz",
		"https://search.yahoo.com/search?p=foo+bar+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLHandleMultipleAndSingleTerms(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/combo_terms.txt")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-s", "foo", "-m"}
	s, err := search.NewSearcher(
		search.WithInput(f),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://www.ask.com/web?q=foo+bar",
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
		"https://www.ask.com/web?q=foo+go+golang",
		"https://bing.com/search?q=foo+go+golang",
		"https://search.brave.com/search?q=foo+go+golang",
		"https://html.duckduckgo.com/html?q=foo+go+golang",
		"https://search.yahoo.com/search?p=foo+go+golang",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar\n")
	args := []string{"-s", "foo", "-e"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://www.ask.com/web?q="foo+bar"`,
		`https://bing.com/search?q="foo+bar"`,
		`https://search.brave.com/search?q="foo+bar"`,
		`https://html.duckduckgo.com/html?q="foo+bar"`,
		`https://search.yahoo.com/search?p="foo+bar"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermArgsExact(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-e", "bar"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://www.ask.com/web?q="foo+bar"`,
		`https://bing.com/search?q="foo+bar"`,
		`https://search.brave.com/search?q="foo+bar"`,
		`https://html.duckduckgo.com/html?q="foo+bar"`,
		`https://search.yahoo.com/search?p="foo+bar"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar baz\n")
	args := []string{"-s", "foo", "-m", "-e"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://www.ask.com/web?q="foo+bar+baz"`,
		`https://bing.com/search?q="foo+bar+baz"`,
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://search.yahoo.com/search?p="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermArgsExact(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-m", "-e", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://www.ask.com/web?q="foo+bar+baz"`,
		`https://bing.com/search?q="foo+bar+baz"`,
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://search.yahoo.com/search?p="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func cmp(t *testing.T, ch <-chan string, want []string) {
	t.Helper()
	var got []string
	for c := range ch {
		got = append(got, c)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

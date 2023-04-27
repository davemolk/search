package search_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/davemolk/search"
)

//////////////
/* no terms */
//////////////
func TestFormatURLNoTerms(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-n"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://search.brave.com/search?q=foo",
		"https://html.duckduckgo.com/html?q=foo",
		"https://www.mojeek.com/search?q=foo",
		"https://lite.qwant.com/?q=foo",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLNoTermsNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-n", "-p=false"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://bing.com/search?q=foo",
		"https://search.brave.com/search?q=foo",
		"https://html.duckduckgo.com/html?q=foo",
		"https://search.yahoo.com/search?p=foo",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLNoTermsMultiTermBase(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo bar baz", "-n"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://search.brave.com/search?q=foo+bar+baz",
		"https://html.duckduckgo.com/html?q=foo+bar+baz",
		"https://www.mojeek.com/search?q=foo+bar+baz",
		"https://lite.qwant.com/?q=foo+bar+baz",
	}
	cmp(t, s.FormatURL(), want)
}

///////////
/* terms */
///////////
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
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://www.mojeek.com/search?q=foo+bar",
		"https://lite.qwant.com/?q=foo+bar",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermNoPrivacy(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar")
	args := []string{"-s", "foo", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://www.mojeek.com/search?q=foo+bar",
		"https://lite.qwant.com/?q=foo+bar",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermArgNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-p=false", "bar"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://www.mojeek.com/search?q=foo+bar",
		"https://lite.qwant.com/?q=foo+bar",
		"https://search.brave.com/search?q=foo+baz",
		"https://html.duckduckgo.com/html?q=foo+baz",
		"https://www.mojeek.com/search?q=foo+baz",
		"https://lite.qwant.com/?q=foo+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermsNoPrivacy(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar baz")
	args := []string{"-s", "foo", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
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
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://www.mojeek.com/search?q=foo+bar",
		"https://lite.qwant.com/?q=foo+bar",
		"https://search.brave.com/search?q=foo+baz",
		"https://html.duckduckgo.com/html?q=foo+baz",
		"https://www.mojeek.com/search?q=foo+baz",
		"https://lite.qwant.com/?q=foo+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermArgsNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-p=false", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
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
		"https://search.brave.com/search?q=foo+bar+baz",
		"https://html.duckduckgo.com/html?q=foo+bar+baz",
		"https://www.mojeek.com/search?q=foo+bar+baz",
		"https://lite.qwant.com/?q=foo+bar+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLCombineTermsNoPrivacy(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/multi_terms.txt")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-s", "foo", "-m", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(f),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		"https://search.brave.com/search?q=foo+bar+baz",
		"https://html.duckduckgo.com/html?q=foo+bar+baz",
		"https://www.mojeek.com/search?q=foo+bar+baz",
		"https://lite.qwant.com/?q=foo+bar+baz",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLCombineTermsArgsNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-m", "-p=false", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://www.mojeek.com/search?q=foo+bar",
		"https://lite.qwant.com/?q=foo+bar",
		"https://search.brave.com/search?q=foo+go+golang",
		"https://html.duckduckgo.com/html?q=foo+go+golang",
		"https://www.mojeek.com/search?q=foo+go+golang",
		"https://lite.qwant.com/?q=foo+go+golang",
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLHandleMultipleAndSingleTermsNoPrivacy(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/combo_terms.txt")
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-s", "foo", "-m", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(f),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		"https://bing.com/search?q=foo+bar",
		"https://search.brave.com/search?q=foo+bar",
		"https://html.duckduckgo.com/html?q=foo+bar",
		"https://search.yahoo.com/search?p=foo+bar",
		"https://bing.com/search?q=foo+go+golang",
		"https://search.brave.com/search?q=foo+go+golang",
		"https://html.duckduckgo.com/html?q=foo+go+golang",
		"https://search.yahoo.com/search?p=foo+go+golang",
	}
	cmp(t, s.FormatURL(), want)
}

////////////////////////////////////////
/* Exact, SearchExact, and TermsExact */
////////////////////////////////////////
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
		`https://search.brave.com/search?q="foo+bar"`,
		`https://html.duckduckgo.com/html?q="foo+bar"`,
		`https://www.mojeek.com/search?q="foo+bar"`,
		`https://lite.qwant.com/?q="foo+bar"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermExactNoPrivacy(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar\n")
	args := []string{"-s", "foo", "-e", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		`https://search.brave.com/search?q="foo+bar"`,
		`https://html.duckduckgo.com/html?q="foo+bar"`,
		`https://www.mojeek.com/search?q="foo+bar"`,
		`https://lite.qwant.com/?q="foo+bar"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermArgsExactNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-e", "-p=false", "bar"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://www.mojeek.com/search?q="foo+bar+baz"`,
		`https://lite.qwant.com/?q="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermExactNoPrivacy(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar baz\n")
	args := []string{"-s", "foo", "-m", "-e", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
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
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://www.mojeek.com/search?q="foo+bar+baz"`,
		`https://lite.qwant.com/?q="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermArgsExactNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-m", "-e", "-p=false", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://bing.com/search?q="foo+bar+baz"`,
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://search.yahoo.com/search?p="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

/* searchExact */
func TestFormatURLSingleTermSearchExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar\n")
	args := []string{"-s", "foo", "-se"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q="foo"+bar`,
		`https://html.duckduckgo.com/html?q="foo"+bar`,
		`https://www.mojeek.com/search?q="foo"+bar`,
		`https://lite.qwant.com/?q="foo"+bar`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermSearchExactNoPrivacy(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar\n")
	args := []string{"-s", "foo", "-se", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://bing.com/search?q="foo"+bar`,
		`https://search.brave.com/search?q="foo"+bar`,
		`https://html.duckduckgo.com/html?q="foo"+bar`,
		`https://search.yahoo.com/search?p="foo"+bar`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermArgsSearchExact(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-se", "bar"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q="foo"+bar`,
		`https://html.duckduckgo.com/html?q="foo"+bar`,
		`https://www.mojeek.com/search?q="foo"+bar`,
		`https://lite.qwant.com/?q="foo"+bar`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSingleTermArgsSearchExactNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-se", "-p=false", "bar"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://bing.com/search?q="foo"+bar`,
		`https://search.brave.com/search?q="foo"+bar`,
		`https://html.duckduckgo.com/html?q="foo"+bar`,
		`https://search.yahoo.com/search?p="foo"+bar`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleBaseSearchTermSearchExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("baz\n")
	args := []string{"-s", "foo bar", "-se"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q="foo+bar"+baz`,
		`https://html.duckduckgo.com/html?q="foo+bar"+baz`,
		`https://www.mojeek.com/search?q="foo+bar"+baz`,
		`https://lite.qwant.com/?q="foo+bar"+baz`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleBaseSearchTermArgsSearchExactNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo bar", "-se", "-p=false", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://bing.com/search?q="foo+bar"+baz`,
		`https://search.brave.com/search?q="foo+bar"+baz`,
		`https://html.duckduckgo.com/html?q="foo+bar"+baz`,
		`https://search.yahoo.com/search?p="foo+bar"+baz`,
	}
	cmp(t, s.FormatURL(), want)
}

/* multiExact */
func TestFormatURLMultipleTermMultiExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar baz\n")
	args := []string{"-s", "foo", "-m", "-me"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q=foo+"bar+baz"`,
		`https://html.duckduckgo.com/html?q=foo+"bar+baz"`,
		`https://www.mojeek.com/search?q=foo+"bar+baz"`,
		`https://lite.qwant.com/?q=foo+"bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermMultiExactNoPrivacy(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("bar baz\n")
	args := []string{"-s", "foo", "-m", "-me", "-p=false"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://bing.com/search?q=foo+"bar+baz"`,
		`https://search.brave.com/search?q=foo+"bar+baz"`,
		`https://html.duckduckgo.com/html?q=foo+"bar+baz"`,
		`https://search.yahoo.com/search?p=foo+"bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermArgsMultiExact(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-m", "-me", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q=foo+"bar+baz"`,
		`https://html.duckduckgo.com/html?q=foo+"bar+baz"`,
		`https://www.mojeek.com/search?q=foo+"bar+baz"`,
		`https://lite.qwant.com/?q=foo+"bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLMultipleTermArgsMultiExactNoPrivacy(t *testing.T) {
	t.Parallel()
	args := []string{"-s", "foo", "-m", "-me", "-p=false", "bar", "baz"}
	s, err := search.NewSearcher(
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://bing.com/search?q=foo+"bar+baz"`,
		`https://search.brave.com/search?q=foo+"bar+baz"`,
		`https://html.duckduckgo.com/html?q=foo+"bar+baz"`,
		`https://search.yahoo.com/search?p=foo+"bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

/* exact > searchExact > multiExact */
func TestFormatURLExactHasPriorityOverSearchExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("baz\n")
	args := []string{"-s", "foo bar", "-se", "-e"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://www.mojeek.com/search?q="foo+bar+baz"`,
		`https://lite.qwant.com/?q="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLExactHasPriorityOverMultiExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("baz\n")
	args := []string{"-s", "foo bar", "-me", "-e"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q="foo+bar+baz"`,
		`https://html.duckduckgo.com/html?q="foo+bar+baz"`,
		`https://www.mojeek.com/search?q="foo+bar+baz"`,
		`https://lite.qwant.com/?q="foo+bar+baz"`,
	}
	cmp(t, s.FormatURL(), want)
}

func TestFormatURLSearchExactHasPriorityOverMultiExact(t *testing.T) {
	t.Parallel()
	bufInput := bytes.NewBufferString("baz\n")
	args := []string{"-s", "foo bar", "-me", "-se"}
	s, err := search.NewSearcher(
		search.WithInput(bufInput),
		search.FromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.CreateQueries()
	want := []string{
		`https://search.brave.com/search?q="foo+bar"+baz`,
		`https://html.duckduckgo.com/html?q="foo+bar"+baz`,
		`https://www.mojeek.com/search?q="foo+bar"+baz`,
		`https://lite.qwant.com/?q="foo+bar"+baz`,
	}
	cmp(t, s.FormatURL(), want)
}

/* helper */
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

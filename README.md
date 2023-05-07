# search
Use privacy mode (on by default) to search brave, duck duck go, mojeek, and qwant and non-privacy mode to search bing, brave, duck duck go, and yahoo. Prints search result URLs and blurbs to stdout. 

## installation
`go install github.com/davemolk/search@latest`
or run one of the binaries included in the repo

# examples
combine additional search terms with your base query (terms are added one at a time)
`search -s golang cloud cli gophers`
```
https://search.brave.com/search?q=golang+cloud
https://html.duckduckgo.com/html?q=golang+cloud
https://www.mojeek.com/search?q=golang+cloud
https://lite.qwant.com/?q=golang+cloud
https://search.brave.com/search?q=golang+cli
https://html.duckduckgo.com/html?q=golang+cli
https://www.mojeek.com/search?q=golang+cli
https://lite.qwant.com/?q=golang+cli
https://search.brave.com/search?q=golang+gophers
https://html.duckduckgo.com/html?q=golang+gophers
https://www.mojeek.com/search?q=golang+gophers
https://lite.qwant.com/?q=golang+gophers
```
combine additional search terms with your base query (use -m to handle multiple terms)
`search -s golang -m cloud cli gophers `
```
https://search.brave.com/search?q=golang+cloud+cli+gophers
https://html.duckduckgo.com/html?q=golang+cloud+cli+gophers
https://www.mojeek.com/search?q=golang+cloud+cli+gophers
https://lite.qwant.com/?q=golang+cloud+cli+gophers
```
combine a mix of single and multiple additional search terms with your base query (use -m and cat in your list of terms)
```
$ cat terms.txt
microservices
cloud technology
machine learning
cli
```
`cat terms.txt | search -s golang -m`
```
https://search.brave.com/search?q=golang+microservices
https://html.duckduckgo.com/html?q=golang+microservices
https://www.mojeek.com/search?q=golang+microservices
https://lite.qwant.com/?q=golang+microservices
https://search.brave.com/search?q=golang+cloud+technology
https://html.duckduckgo.com/html?q=golang+cloud+technology
https://www.mojeek.com/search?q=golang+cloud+technology
https://lite.qwant.com/?q=golang+cloud+technology
https://search.brave.com/search?q=golang+machine+learning
https://html.duckduckgo.com/html?q=golang+machine+learning
https://www.mojeek.com/search?q=golang+machine+learning
https://lite.qwant.com/?q=golang+machine+learning
https://search.brave.com/search?q=golang+cli
https://html.duckduckgo.com/html?q=golang+cli
https://www.mojeek.com/search?q=golang+cli
https://lite.qwant.com/?q=golang+cli
```

## flags
```
[customize basic query info]
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
    search -s "foo bar" baz => https://seach.brave.com/search?q=foo+bar+baz


[customize exact searching]
-e  exact searching for entire query
	search -s foo bar -e => https://search.brave.com/search?q="foo+bar", etc.
	default: false

-me exact matching for additional terms
	search -s foo -me bar baz => https://search.brave.com/search?q=foo+"bar+baz", etc.
    default: false

-se exact matching for search term(s)
	search -s "foo bar" -se baz => https://search.brave.com/search?q="foo+bar"+baz, etc.
	default: false


[customize requests]
-c  max number of concurrent requests
	default: 10

-os operating system (used for creating browser headers)
	arguments: any, l, m, or w
	default: w

-t  request timeout, in ms
	default: 5000


[customize output]
-l  length of result summary
	default: 500

-u  include result urls in output
	default: true


help
-d  print the search url to help debug queries
	default: false
-h  help`)
```

## Note
Each request gets a randomly assigned user agent corresponding to your os as well as appropriate headers (50/50 chance of chrome or firefox, thanks [fuzzyHelpers](https://github.com/davemolk/fuzzyHelpers)). Go unfortunately doesn't preserve header order, so if that's important to you and what you're up to, you'll have to look elsewhere.
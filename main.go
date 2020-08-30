package main

import (
	"flag"
	"fmt"
)

func check(e error) {
	if e != nil {
		// fmt.Fprintf(os.Stderr, "err: %s\n", err)
		// return
		panic(e)
	}
}

func main() {
	var rootsOnly bool
	flag.BoolVar(&rootsOnly, "roots", false, "Only show canonical root domains")
	flag.Parse()

	dirs := []dirFn{
		dirDod,
		dirAf,
		dirArmy,
		dirNavy,
	}

	var err error
	urls := make(chan string)

	err = scrapeDirs(dirs, urls)
	check(err)

	seenDomain := make(map[string]bool)
	seenRoot:= make(map[string]bool)

	// process and output results from directory scraper workers
	for u := range urls {
		// parse result
	    sub := getSub(u)

		// de-duplicate results
		if _, dd := seenDomain[sub]; dd {
			continue
		}
		seenDomain[sub] = true

		// discard non-dotmil domains
		if !isDotmil(sub) {
			continue
		}

		root := extractRoot(sub)

		// check if duplicate root
		_, dr := seenRoot[root]

		// output root only if unique and roots flag is set
		// otherwise, output full domain
		if rootsOnly {
		    if !dr {
				fmt.Println(root)
			}
		} else {
			fmt.Println(sub)
		}

		seenRoot[root] = true
	}
}

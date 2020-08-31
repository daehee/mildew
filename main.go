package main

import (
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
	var err error

	// DoD directory scraping functions to loop through
	dirs := []dirFn{
		dirDod,
		dirAf,
		dirArmy,
		dirNavy,
	}

	// send results from directory scraper to urls channel
	urls := make(chan string)
	// start goroutines for each directory scraping function
	err = scrapeDirs(dirs, urls)
	check(err)

	// save subs and roots in map as a reference to de-dupe data
	seenDomain := make(map[string]bool)
	seenRoot:= make(map[string]bool)

	// set up channel to receive roots, as jobs for requesting crt data
	crtJobs := make(chan string)
	// set up channel to receive results from crt scraper
	crtOut := make(chan string)
	// start goroutine for crt scraper, and
	// prepare crtJobs channel to receive roots as jobs
	go scrapeCrts(crtJobs, crtOut)

	// read directory scraper results from urls channel
	for u := range urls {
	    sub := getSub(u)

		// discard non-dotmil domains
		if !isDotmil(sub) {
			continue
		}

		// skip loop if sub is duplicate
		if _, dd := seenDomain[sub]; dd {
			continue
		}
		seenDomain[sub] = true

		// output the sub before moving on to root processing and crt scraping
		fmt.Println(sub)

		root := extractRoot(sub)
		// de-dupe root to make sure request crts for a root only once
		_, dr := seenRoot[root]
		if !dr {
			// send root to crt scraper through crtJobs channel
			crtJobs <- root
		}
		seenRoot[root] = true
	}
	// all jobs sent, close channel
	close(crtJobs)

	for r := range crtOut {

		// de-dupe results from crt scraping
		if _, dd := seenDomain[r]; dd {
			continue
		}
		seenDomain[r] = true

		fmt.Println(r)
	}
}

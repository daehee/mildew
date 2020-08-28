package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"net/url"
)

func check(e error) {
	if e != nil {
		// fmt.Fprintf(os.Stderr, "err: %s\n", err)
		// return
		panic(e)
	}
}

func getSub(u string) string {
	p, err := url.Parse(u)
	check(err)

	sub := cleanDomain(p.Hostname())
	return sub
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)
	return d
}

func isDotmil(d string) bool {
	return strings.HasSuffix(d, "mil")
}

func extractRoot(d string) string {
	split := strings.Split(d, ".")
	split = split[len(split)-2:]
	root := strings.Join(split, ".")
	return root
}

type dirFn func(chan<- string) error

func scrapeDirs(dirs []dirFn) <- chan string {
	res := make(chan string)
	var wg sync.WaitGroup

	for _, dir := range dirs {
		wg.Add(1)
		fn := dir
		go func() {
			defer wg.Done()
			err := fn(res)
			check(err)
		}()
	}

	// The dir functions have returned, so all calls to wg.Add are done. Start a
	// goroutine to close res once all the sends are done
	go func() {
		wg.Wait()
		close(res)
	}()

	return res
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

	urls := scrapeDirs(dirs)

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

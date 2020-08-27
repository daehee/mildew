package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"net/url"

	"github.com/gocolly/colly"
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

type getFn func(*colly.Collector, chan<- string) error

func main() {
	var rootsOnly bool
	flag.BoolVar(&rootsOnly, "roots", false, "Only show canonical root domains")
	flag.Parse()

	sources := []getFn{
		getDod,
		getAf,
		getArmy,
		getNavy,
	}

	out := make(chan string)
	var wg sync.WaitGroup
	c := colly.NewCollector()

	for _, source := range sources {
		wg.Add(1)
		fn := source
		go func() {
			defer wg.Done()
			err := fn(c, out)
			check(err)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	seen := make(map[string]bool)
	for d := range out {
	    d = getSub(d)

	    // extract roots only if flag is set by user
		if rootsOnly {
			d = extractRoot(d)
		}

		// de-duplicate results
		if _, ok := seen[d]; ok {
			continue
		}
		seen[d] = true

		// discard non-dotmil domains
		if !isDotmil(d) {
			continue
		}

		fmt.Println(d)
	}

}

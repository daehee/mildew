package main

import (
	"context"
	"log"

	"github.com/daehee/mildew/mildew"
)

// mildew scrapes domains from official DoD website directories
// and certificate transparency logs
//
// 1: Scrape DoD directories for subdomains
// 2: Request additional certificate transparency subdomains using root domains from 1
// 3: Output to stdout and file
func main() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mw := mildew.NewMildew()

	log.Printf("scraping DoD web directories")
	err = mw.ScrapeDirs(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("scraping certificate transparency data")
	err = mw.ScrapeCrts(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// mw.OutputScreen()
	mw.OutputFile("mildew.out")
}

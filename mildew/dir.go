package mildew

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ScrapeDirs scrapes all DoD website directories and saves to Mildew object's Subs field
func (mw *Mildew) ScrapeDirs(ctx context.Context) error {
	dirStream := make(chan string)
	// Launch goroutine for scraping directories
	go func() {
		defer close(dirStream)
		// Initialize base colly collector to be used by each directory scraper function
		// TODO tune colly options
		c := colly.NewCollector()
		var err error

		err = dirDod(c, dirStream)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = dirAf(c, dirStream)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = dirArmy(c, dirStream)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = dirNavy(c, dirStream)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// Process incoming URLs from web directories
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("scrape cancelled")
		case v, ok := <-dirStream:
			// nil value signals closed channels, so return from function
			if ok == false {
				return nil
			}
			sub := urlToSub(v)
			if !isDotmil(sub) {
				continue
			}
			if mw.Subs.Has(sub) {
				continue
			}
			mw.Subs.Insert(sub)
		}
	}
}

// dirDod scrapes defense.gov directory
func dirDod(c *colly.Collector, dirStream chan<- string) error {
	cc := c.Clone()
	cc.OnHTML("div.DGOVLinkBox > div", func(e *colly.HTMLElement) {
		dirStream <- e.ChildAttr("a[href]", "href")
	})
	err := cc.Visit("https://www.defense.gov/Resources/Military-Departments/DOD-Websites/")
	if err != nil {
		return fmt.Errorf("error scraping defense.gov: %v", err)
	}
	return nil
}

// dirDod scrapes af.mil directory
func dirAf(c *colly.Collector, dirStream chan<- string) error {
	cc := c.Clone()
	cc.OnHTML("a.AFSiteLink, a.AFSiteBaseLink", func(e *colly.HTMLElement) {
		dirStream <- e.Attr("href")
	})
	// These links are duplicates of each other, go straight to srBaseList
	// cc.OnHTML("a.AFAlphaLink", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })
	err := cc.Visit("http://www.af.mil/AFSites.aspx")
	if err != nil {
		return fmt.Errorf("error scraping af.mil: %v", err)
	}
	err = cc.Visit("https://www.af.mil/AF-Sites/srBaseList/A/#A")
	if err != nil {
		return fmt.Errorf("error scraping af.mil: %v", err)
	}
	return err
}

// dirArmy scrapes army.mil directory
func dirArmy(c *colly.Collector, dirStream chan<- string) error {
	cc := c.Clone()
	cc.OnHTML("div.links-list a", func(e *colly.HTMLElement) {
		dirStream <- e.Attr("href")
	})
	err := cc.Visit("http://www.army.mil/info/a-z/")
	if err != nil {
		return fmt.Errorf("error scraping army.mil: %v", err)
	}
	return nil
}

// dirNavy scrapes website URLs from Navy's VueJS SPA,
// requires some JSON parsing
func dirNavy(c *colly.Collector, dirStream chan<- string) error {
	cc := c.Clone()
	cc.OnHTML("#dnn_ctr752_ModuleContent > script:nth-of-type(2)", func(e *colly.HTMLElement) {
		// find JSON string feeding VueJS website directory
		jsonData := e.Text[strings.Index(e.Text, "[{") : strings.Index(e.Text, "}]}]")+4]

		// parse JSON
		var data []struct {
			SiteUrl string `json:"url"`
		}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			return
		}

		for _, d := range data {
			// discard empty site URL fields, since some Navy units only list their social media profiles
			if len(d.SiteUrl) > 0 {
				dirStream <- d.SiteUrl
			}
		}
	})
	err := cc.Visit("https://www.navy.mil/Resources/Navy-Directory/")
	if err != nil {
		return fmt.Errorf("error scraping navy.mil: %v", err)
	}
	return nil
}

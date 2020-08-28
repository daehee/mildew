package main

import (
    "encoding/json"
    "log"
    "strings"

    "github.com/gocolly/colly"
)

func dirDod(out chan<- string) error {
    c := colly.NewCollector()

    c.OnHTML("div.DGOVLinkBox > div", func(e *colly.HTMLElement) {
        link := e.ChildAttr("a[href]", "href")
        out <- link
    })

    err := c.Visit("https://www.defense.gov/Resources/Military-Departments/DOD-Websites/")

    return err
}

func dirAf(out chan<- string) error {
    c := colly.NewCollector()

    c.OnHTML("a.AFSiteLink, a.AFSiteBaseLink", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        out <- link
    })

    c.OnHTML("a.AFAlphaLink", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    err := c.Visit("http://www.af.mil/AFSites.aspx")

    return err
}

func dirArmy(out chan<- string) error {
    c := colly.NewCollector()

    c.OnHTML("div.links-list a", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        out <- link
    })

    err := c.Visit("http://www.army.mil/info/a-z/")

    return err
}

// Scrapes website URLs from Navy's VueJS SPA, which requires some messy JSON parsing
func dirNavy(out chan<- string) error {
    c := colly.NewCollector()

    c.OnHTML("#dnn_ctr752_ModuleContent > script:nth-of-type(2)", func(e *colly.HTMLElement) {
        // find JSON string feeding VueJS website directory
        jsonData := e.Text[strings.Index(e.Text, "[{") : strings.Index(e.Text, "}]}]")+4]

        // parse JSON
        var data []struct {
            SiteUrl string `json:"url"`
        }
        err := json.Unmarshal([]byte(jsonData), &data)
        if err != nil {
            log.Fatal(err)
        }

        for _, d := range data {
            // discard empty site URL fields, since some Navy units only list their social media profiles
            if len(d.SiteUrl)> 0 {
                out <- d.SiteUrl
            }
        }
    })

    err := c.Visit("https://www.navy.mil/Resources/Navy-Directory/")

    return err
}

package main

import (
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
    "time"

    "github.com/gocolly/colly"
    "github.com/gocolly/colly/queue"
)

func scrapeCrts(jobs <-chan string, res chan<- string) {
    c := colly.NewCollector(
        // attach a debugger to the collector
        // colly.Debugger(&debug.LogDebugger{}),
    )

    // rate limit colly and set delay
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*crt.*",
        Delay:      5 * time.Second,
    })

    // create a request queue with 2 consumer threads
    q, _ := queue.New(
        2,  // number of consumer threads
        &queue.InMemoryQueueStorage{MaxSize: 10000},    // use default queue storage
    )

    // c.OnRequest(func(r *colly.Request) {
    //     fmt.Println("Visiting", r.URL)
    // })

    c.OnResponse(func(r *colly.Response){
        // Parse subdomains from JSON data
        jsonData := string(r.Body)
        var data[]struct {
            CaId int `json:"issuer_ca_id"`
            NameValue string `json:"name_value"`
        }
        err := json.Unmarshal([]byte(jsonData), &data)
        check(err)

        for _, d := range data {
            // account for name_values data containing newlines
            split := strings.Split(d.NameValue, "\n")
            for _, s := range split {
                // regex extract dotmil domain only, case-insensitive
                re := regexp.MustCompile(`(?i)((?:([a-z0-9]\.|[a-z0-9][a-z0-9\-]{0,61}[a-z0-9])\.)+)(mil)\.?`)
                match := re.FindStringSubmatch(s)
                s = match[0]
                s = cleanDomain(s)
                res <- s
            }
        }

    })

    // receive roots as crt scraper jobs
    for j := range jobs {
        root := j
        // queue colly request to download JSON format for the root from crt.sh
        crtUrl := fmt.Sprintf("https://crt.sh/?dNSName=%%25.%s&output=json", root)
        q.AddURL(crtUrl)
    }

    // Execute colly queue
    q.Run(c)

    // Wait until threads are finished; may be redundant with Queue Run, which
    // blocks while queue has active requests
    // c.Wait()

    // Done sending crt scraping results, close channel
    close(res)
}
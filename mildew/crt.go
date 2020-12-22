package mildew

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

func (mw *Mildew) ScrapeCrts(ctx context.Context) error {
	c := colly.NewCollector(colly.StdlibContext(ctx))

	// rate limit colly and set delay
	c.Limit(&colly.LimitRule{
		DomainGlob: "*crt.*",
		Delay:      5 * time.Second,
	})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		log.Printf("%s", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		// Parse subdomains from JSON data
		jsonData := string(r.Body)
		var data []struct {
			CaId      int    `json:"issuer_ca_id"`
			NameValue string `json:"name_value"`
		}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			return
		}

		for _, d := range data {
			// account for name_values data containing newlines
			split := strings.Split(d.NameValue, "\n")
			for _, s := range split {
				match := dotmilRx.FindStringSubmatch(s)
				s = match[0]
				sub := cleanDomain(s)
				if mw.Subs.Has(sub) {
					continue
				}
				mw.Subs.Insert(sub)
			}
		}

	})

	// receive roots as crt scraper rootStream
	for _, root := range mw.GetRoots() {
		// queue colly request to download JSON format for the root from crt.sh
		err := q.AddURL(fmt.Sprintf("https://crt.sh/?dNSName=%%25.%s&output=json", root))
		if err != nil {
			return fmt.Errorf("error queueing crt scrape for %s: %v", root, err)
		}
	}
	// Execute colly queue
	err := q.Run(c)
	if err != nil {
		return fmt.Errorf("error executing crt scrape: %v", err)
	}

	// go func() {
	// 	// Wait until threads are finished; may be redundant with Queue Run, which
	// 	// blocks while queue has active requests
	// 	c.Wait()
	// 	// Done sending crt scraping results, close channel
	// 	done <- struct{}{}
	// }()

	return nil
}

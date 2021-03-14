package main

import (
	"fmt"
	"time"

	"github.com/bgroupe/obey/services"
)

type scraper struct {
	Ticker *time.Ticker
	Done   chan bool
	cfg    tomlConfig
	Client services.Client
}

func newScraper() (*scraper, error) {
	duration, err := time.ParseDuration(config.ScrapeConfig.Interval)
	if err != nil {
		return nil, err
	}
	// TODO: configure client based on config
	cli, err := services.NewDockerClient()

	if err != nil {
		return nil, err
	}
	sc := scraper{
		Ticker: time.NewTicker(duration),
		Done:   make(chan bool),
		cfg:    config,
		Client: cli,
	}

	return &sc, err
}

func (sc *scraper) Run() {
	for {
		select {
		case <-sc.Done:
			return
		case tick := <-sc.Ticker.C:
			fmt.Println("scraping services@:", tick)
			scraped, err := sc.Client.Scrape()
			if err != nil {
				panic(err)
			}
			// for _, container := range scraped {
			// 	fmt.Println(container.Name)
			// }
			// Serialize Scraped Stuff
			// Send to Scheduler via GRPC
			reportServiceData(scraped)
		}
	}
}

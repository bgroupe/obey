package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/KromDaniel/rejonson"
	"github.com/gadabout/obey/config"
)

type pollWorker struct {
	Ticker *time.Ticker
	Done   chan bool
	cfg    config.ServiceConfig
}

func NewPollWorker(interval string, cfg config.ServiceConfig) *pollWorker {
	duration, _ := time.ParseDuration(interval)
	pw := &pollWorker{
		Ticker: time.NewTicker(duration),
		Done:   make(chan bool),
		cfg:    cfg,
	}

	return pw
}

func (pw *pollWorker) Run(client rejonson.Client, manifestKey string) {
	for {
		select {
		case <-pw.Done:
			return
		case tick := <-pw.Ticker.C:
			fmt.Println("fetching services", tick)
			environments := FetchServices(pw.cfg)
			jsonString, err := json.Marshal(environments)
			if err != nil {
				fmt.Println("unable to set marshal json", err)
			}
			client.JsonSet(manifestKey, ".", string(jsonString))
		}
	}
}

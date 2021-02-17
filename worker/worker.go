package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	// workerID is the id assigned by the scheduler
	// after registering on scheduler.
	workerID string
)

func init() {
	loadConfig()
}

// Entry point of the worker application.
func main() {
	scraper, err := newScraper()
	if err != nil {
		fatal(fmt.Sprintf("Error starting scraper: (%e)\n", err))
	}

	go startGRPCServer()
	go registerWorker()
	go scraper.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case s := <-sig:
			fatal(fmt.Sprintf("Signal (%d) received, stopping\n", s))
		}
	}
}

func fatal(message string) {
	deregisterWorker()
	log.Fatalln(message)
}

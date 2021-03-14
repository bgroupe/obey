package main

import (
	"github.com/bgroupe/obey/services"
	log "github.com/sirupsen/logrus"
)

type listener struct {
	cfg    tomlConfig
	Client services.Client
}

func newListener() (*listener, error) {
	cli, err := services.NewDockerClient()
	if err != nil {
		return nil, err
	}

	l := listener{
		cfg:    config,
		Client: cli,
	}

	return &l, err
}

func (l *listener) Listen() {
	// primitive
	// TODO: implement retry w/ backoff
	// https://github.com/avast/retry-go
	log.WithFields(log.Fields{
		"process": "Listener",
		"info":    "startup",
	}).Info("performing initial discovery")

	// Discover initial set of services before polling for changes
	discovered, err := l.Client.Discover()

	if err != nil {
		log.WithFields(log.Fields{
			"process": "Listener",
			"status":  "error",
		}).Error(err)
	}

	reportServiceData(discovered)

	log.WithFields(log.Fields{
		"process": "Listener",
		"info":    "startup",
	}).Info("initiating listener...")

	l.Client.Poll()
}

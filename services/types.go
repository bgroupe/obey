package services

import (
	"time"

	"github.com/docker/docker/client"
)

// Client - Main interface for worker
type Client interface {
	Poll(labels chan string)
	Batch(labels chan string, maxItems int, maxTimeout time.Duration) chan []ScrapedContainer
	Discover() ([]ScrapedContainer, error)
	Scrape() ([]ScrapedContainer, error)
}

// DockerClientConfig - config object for Docker client
type DockerClientConfig struct {
	scrape                   bool
	scrapeAll                bool
	poll                     bool
	deriveMetadataFromLabels bool
}

// DockerClient - struct for hold client config, implements Client interface
type DockerClient struct {
	Client *client.Client
	Config DockerClientConfig
}

// ScrapedContainer used for serializing container data back to the worker
//  TODO: use the proto struct instead
type ScrapedContainer struct {
	Name     string
	Version  string
	State    string
	Status   string
	Created  string
	ImageSha string
}

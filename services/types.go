package services

import (
	"github.com/docker/docker/client"
)

// Client - Main interface for worker
type Client interface {
	Scrape() ([]ScrapedContainer, error)
}

// DockerClientConfig - config object for Docker client
type DockerClientConfig struct {
	scrapeAll                bool
	deriveMetadataFromLabels bool
}

// DockerClient - struct for hold client config, implements Client interface
type DockerClient struct {
	Client *client.Client
	Config DockerClientConfig
}

// ScrapedContainer used for serializing container data back to the worker
type ScrapedContainer struct {
	Name     string
	Version  string
	State    string
	Status   string
	Created  string
	ImageSha string
}

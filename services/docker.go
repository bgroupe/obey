package services

// https://pkg.go.dev/github.com/docker/docker
import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	scrapeLabel  = "obey.com/scrape"
	versionLabel = "obey.com/version"
	serviceLabel = "obey.com/serviceName"
)

// Scrape scrapes all containers based on configured labels
func (d *DockerClient) Scrape() ([]ScrapedContainer, error) {
	ctx := context.Background()
	containers, err := d.Client.ContainerList(ctx, types.ContainerListOptions{})

	if err != nil {
		return []ScrapedContainer{}, err
	}

	// TODO: Export scraped containers to filter function for scrapeAll logic
	var scSlice []ScrapedContainer
	for _, container := range containers {
		labels := container.Labels
		if val, _ := labels[scrapeLabel]; val == "true" {
			var sc ScrapedContainer
			shortSha := strings.Split(container.ImageID, "sha256:")[1]
			sc.Name = d.deriveServiceName(container)
			sc.Version = d.deriveServiceVersion(container)
			sc.State = container.State
			sc.Status = container.Status
			sc.ImageSha = shortSha[:7]
			scSlice = append(scSlice, sc)
		}
	}

	return scSlice, err
}

func (d *DockerClient) deriveServiceVersion(container types.Container) string {
	if d.Config.deriveMetadataFromLabels {
		if version, ok := container.Labels[versionLabel]; ok {
			return version
		}
		return "Label `version` Not Found"
	}
	tag := strings.SplitAfter(container.Image, ":")[1]
	return tag
}

func (d *DockerClient) deriveServiceName(container types.Container) string {
	if d.Config.deriveMetadataFromLabels {
		if name, ok := container.Labels[serviceLabel]; ok {
			return name
		}
		return "Label `serviceName` Not Found"
	}
	// taking  the first name here
	return container.Names[0]
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// NewDockerClient returns a new DockerClient
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	dc := DockerClient{
		Client: cli,
		Config: DockerClientConfig{
			scrapeAll:                false,
			deriveMetadataFromLabels: true,
		},
	}

	return &dc, nil

}

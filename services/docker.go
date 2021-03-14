package services

// https://pkg.go.dev/github.com/docker/docker
import (
	"context"
	"strings"
	"time"

	pb "github.com/bgroupe/obey/jobscheduler"
	"google.golang.org/grpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// TODO: Use Labels from config
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

// ScrapeWithLabel scrapes a container based on label inputs
func (d *DockerClient) ScrapeWithLabel(label string) (ScrapedContainer, error) {
	ctx := context.Background()
	var sc ScrapedContainer
	containers, err := d.Client.ContainerList(ctx, types.ContainerListOptions{})

	if err != nil {
		return sc, err
	}

	for _, container := range containers {
		labels := container.Labels
		if val, _ := labels[label]; val == "true" {

			shortSha := strings.Split(container.ImageID, "sha256:")[1]
			sc.Name = d.deriveServiceName(container)
			sc.Version = d.deriveServiceVersion(container)
			sc.State = container.State
			sc.Status = container.Status
			sc.ImageSha = shortSha[:7]
		}
	}

	return sc, err
}

// Poll wraps polling functionality
func (d *DockerClient) Poll() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventChan, errChan := d.Client.Events(ctx, types.EventsOptions{})
	log.WithFields(log.Fields{
		"process": "Listener",
		"status":  "connected",
	}).Info("connected to Docker Daemon")
	for {
		select {
		case <-errChan:
			// change this
			panic(<-errChan)
		case event := <-eventChan:
			if event.Type == "container" && event.Action == "create" {
				log.WithFields(log.Fields{
					"process": "listener",
					"status":  event.Status,
					"from":    event.From,
					"type":    event.Type,
					"service": event.Actor.Attributes["obey.com/serviceName"],
					"version": event.Actor.Attributes["obey.com/version"],
					// "attributes": event.Actor.Attributes,
				}).Info("New Service Updated: ")
				//    sc, err := d.ScrapeWithLabel(event.Actor.Attributes["obey.com/serviceName"])
				// need to batch here
			}
		}
	}
}

// Discover performs the initial load of services in an environment
func (d *DockerClient) Discover() ([]ScrapedContainer, error) {
	scraped, err := d.Scrape()
	if err != nil {
		panic(err)
	}
	for _, container := range scraped {
		log.WithFields(log.Fields{
			"process": "Listener",
			"name":    container.Name,
		}).Info("Discovered:")
	}
	return scraped, err
}

// Private

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

func sendDataReport(sc []ScrapedContainer) {

	conn, err := grpc.Dial("127.0.0.1:50000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Need to append pointers here
	var serviceData []*pb.ServiceData

	for _, container := range sc {
		sd := pb.ServiceData{
			Name:     container.Name,
			Version:  container.Version,
			State:    container.State,
			Status:   container.Status,
			Created:  container.Created,
			Revision: container.ImageSha,
		}

		serviceData = append(serviceData, &sd)
	}

	reportReq := pb.ReportServiceDataRequest{
		// Name:        config.WorkerEnvConfig.Name,
		Name:        "local",
		ServiceData: serviceData,
	}

	r, err := c.ReportServiceData(ctx, &reportReq)

	if err != nil {
		log.Fatalf("could not report service data: %v", err)
	}

	log.Printf("Scrape Report OK: %t", r.Success)
}

package models

// Worker hold information about workers
//  `BroadcastAddress` for private workers behind public ingress
//  `Type` distribution type: k8s, docker, manual
//  `Uptime` distribution type: k8s, docker, manual
type Worker struct {
	ID               string `redis:"id" json:"id"`
	Address          string `redis:"address" json:"address"`
	BroadcastAddress string `redis:"broadcastAddress" json:"broadcast-address"`
	Type             string `redis:"type" json:"type"`
	Env              Environment
}

// Environment holds information about a cluster/machine/architecture/topology
type Environment struct {
	Name     string    `json:"name"`
	Services []Service `json:"services"`
}

// Service holds information about a given
type Service struct {
	Name     string         `json:"name"`
	Host     string         `json:"host"`
	Metadata ServiceVersion `json:"version_metadata"`
}

// ServiceVersion holds version metadata of a given service
type ServiceVersion struct {
	Metadata struct {
		Branch string `json:"branch"`
		Commit string `json:"commit"`
	}
	Version string `json:"version"`
}

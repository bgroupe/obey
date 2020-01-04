package services

import (
	"net/http"
	"encoding/json"
	"github.com/gadabout/obey/config"
	"fmt"
)

type Environment struct {
	Name     string    `json:"name"`
	Services []Service `json:"services"`
}

type Service struct {
	Name     string         `json:"name"`
	Host     string         `json:"host"`
	Metadata ServiceVersion `json:"version_metadata"`
}

type ServiceVersion struct {
	metadata struct {
		Branch string `json: "branch"`
		Commit string `json: "commit"`
	}
	Version string `json: "version"`
}

func FetchServices(cfg config.ServiceConfig) []Environment {
	client := &http.Client{}
	var environments []Environment

	for _, env := range cfg.Environments {
		serviceEnv := Environment{Name: env.EnvName}
		for _, service := range env.Services {
			url := fmt.Sprintf("https://%s/version.json", service.Value)
			resp, err := client.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			var sv ServiceVersion

			json.NewDecoder(resp.Body).Decode(&sv)

			serviceBlock := Service{
				Name:     service.Name,
				Host:     service.Value,
				Metadata: sv,
			}
			serviceEnv.Services = append(serviceEnv.Services, serviceBlock)
		}
		environments = append(environments, serviceEnv)
	}
	return environments
}

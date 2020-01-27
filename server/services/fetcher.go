package services

import (
	"net/http"
	"encoding/json"
	"github.com/gadabout/obey/config"
	"fmt"
)

const versionFile = "version.json"
const httpProtocol = "http"

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
	Metadata struct {
		Branch string `json:"branch"`
		Commit string `json:"commit"`
	}
	Version string `json:"version"`
}



func FetchServices(cfg config.ServiceConfig) []Environment {
	ec := make(chan Environment)
	fetchedEnvironments := []Environment{}
	for _, env := range cfg.Environments {
		fmt.Println("fetching environment", env.EnvName)
		go fetchEnvironment(env, ec)
		fetchedEnvironments = append(fetchedEnvironments, <-ec)
	}

	return fetchedEnvironments
}

// private

func fetchEnvironment(env config.ConfigEnv, ec chan Environment) {
	c := make(chan Service)
	fetchedServices := fetchServices(env.Services, c)
	e := Environment{
		Name: env.EnvName,
		Services: fetchedServices,
		// TODO: include more environment metadata from spec
	}
	ec <- e
}

func fetchServices(configServices []config.ConfigSvc, c chan Service) []Service {
	fetchedServices := []Service{}
	for _, svc := range configServices {
	  fmt.Println("checking host:", svc.Value)
	  go getServiceVersion(svc, c)
	  fetchedServices = append(fetchedServices, <-c)
	}

	return fetchedServices
}



func getServiceVersion(configSvc config.ConfigSvc, c chan Service) {
	client := &http.Client{}
	url := fmt.Sprintf("%s://%s/%s", httpProtocol, configSvc.Value, versionFile )
	res, err := client.Get(url)

	if err != nil {
		panic(err)
	} else {
		fmt.Println(url, "checked")
	}

	defer res.Body.Close()
	var sv ServiceVersion

	json.NewDecoder(res.Body).Decode(&sv)

	serviceBlock := Service{
		Name: configSvc.Name,
		Host: configSvc.Value,
		Metadata: sv,
	}

	c <- serviceBlock

}

func deprecatedFetch(cfg config.ServiceConfig) []Environment {
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

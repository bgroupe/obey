package config

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
)

// type ServiceConfig struct {
// 	Environments []struct {
// 		EnvName string `yaml:"name"`
// 		Discovery bool `yaml:"discovery"`
// 		Selector string `yaml:"selector"`
// 		Services []struct {
// 			Name string `yaml:"name"`
// 			Value string `yaml:"value"`
// 		}
// 	}
// }

type ServiceConfig struct {
	Environments []ConfigEnv
}

type ConfigEnv struct {
	EnvName string `yaml:"name"`
	Discovery bool `yaml:"discovery"`
	Selector string `yaml:"selector"`
	Services []ConfigSvc
}

type ConfigSvc struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}


func processError(err error) {
    fmt.Println(err)
    os.Exit(2)
}

func ReadConfig(cfg *ServiceConfig, filePath string) {
	f, err := os.Open(filePath)

	if err != nil {
		processError(err)
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(cfg)

    if err != nil {
		processError(err)
	}
}

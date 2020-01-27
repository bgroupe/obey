package main

import (
	"fmt"
	"github.com/gadabout/obey/config"
	"github.com/gadabout/obey/services"
	"github.com/gin-gonic/gin"
)



func main() {
	// load service configuration from yaml file into struct
	configFile := "example-config.yaml"
	// TODO: command line flags to pass filepath
	var cfg config.ServiceConfig
	config.ReadConfig(&cfg, configFile)

	for _, env := range cfg.Environments {
		fmt.Printf("selector: %+v\n", env.Selector)
	}

	router := gin.Default()
	// TODO: register routes from config


	router.GET("/services", func(c *gin.Context) {

		environments := services.FetchServices(cfg)
		// fmt.Printf("%+v\n", environments)

		// Shouldn't have to marshal here
	    c.JSON(200, environments)

	})
	router.Run()
}

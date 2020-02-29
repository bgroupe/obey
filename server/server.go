package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/KromDaniel/rejonson"
	"github.com/gadabout/obey/config"
	"github.com/gadabout/obey/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// TODO: handle errors with standard or gin logger..
func handleError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

const manifestKey = "obey_manifest"

func main() {
	// load service configuration from yaml file into struct
	configFile := "example-config.yaml"
	// TODO: command line flags to pass filepath
	var cfg config.ServiceConfig
	config.ReadConfig(&cfg, configFile)

	pollWorker := services.NewPollWorker("5m", cfg)

	for _, env := range cfg.Environments {
		fmt.Printf("selector: %+v\n", env.Selector)
	}
	// TODO: load specifics from env
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rClient := rejonson.ExtendClient(redisClient)
	defer rClient.Close()

	// TODO: dump Struct into json

	rClient.JsonSet("testThing", ".", "{\"fart\": {\"foo\": \"bar\"}}")
	thing, err := rClient.JsonGet("testThing", "fart").Result()

	handleError(err)

	fmt.Println(thing)
	rClient.JsonSet("testThing", ".fart.foo", "\"baz\"")
	thing2, err := rClient.JsonGet("testThing", "fart").Result()

	handleError(err)

	fmt.Println(thing2)

	router := gin.Default()

	router.GET("/services", func(c *gin.Context) {

		environments := services.FetchServices(cfg)
		jsonString, err := json.Marshal(environments)
		handleError(err)
		rClient.JsonSet(manifestKey, ".", string(jsonString))
		fmt.Println(jsonString)
		// Shouldn't have to marshal here
		c.JSON(200, environments)

	})

	router.GET("/r/services", func(c *gin.Context) {

		jsonString, err := rClient.JsonGet(manifestKey).Result()
		fmt.Println(jsonString)
		// temporarily dump into array of interface to parse string
		var result []interface{}
		err = json.Unmarshal([]byte(jsonString), &result)
		handleError(err)

		c.JSON(200, result)
	})

	router.GET("/services/:environment", func(c *gin.Context) {
		queriedEnv := c.Param("environment")
		jsonString, err := rClient.JsonGet(manifestKey).Result()
		var environments []services.Environment
		err = json.Unmarshal([]byte(jsonString), &environments)
		handleError(err)

		var envResults services.Environment

		for _, env := range environments {
			if env.Name == queriedEnv {
				envResults = env
				fmt.Println(envResults)
			}
		}
		// TODO: make this smarter
		c.JSON(200, envResults)

	})

	router.GET("/services/:environment/:service", func(c *gin.Context) {
		queriedEnv := c.Param("environment")
		queriedService := c.Param("service")

		jsonString, err := rClient.JsonGet(manifestKey).Result()
		var environments []services.Environment
		err = json.Unmarshal([]byte(jsonString), &environments)
		handleError(err)

		var foundService services.Service

		for _, env := range environments {
			if env.Name == queriedEnv {
				for _, svc := range env.Services {
					if svc.Name == queriedService {
						foundService = svc
						fmt.Println(foundService)
					}
				}
			}

		}
		c.JSON(200, foundService)
	})

	router.GET("/kill-worker", func(c *gin.Context) {
		// Here is an arbitrary function to stop the worker
		pollWorker.Ticker.Stop()
		pollWorker.Done <- true
		fmt.Println("Worker stopped")
		c.String(200, "Worker Stopped")
	})

	go pollWorker.Run(*rClient, manifestKey)

	router.Run()
}

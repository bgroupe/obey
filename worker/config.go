package main

// TODO: https://github.com/pelletier/go-toml

import (
	"log"

	"github.com/BurntSushi/toml"
)

type grpcServerConfig struct {
	Addr    string `toml:"addr"`
	UseTLS  bool   `toml:"use_tls"`
	CrtFile string `toml:"crt_file"`
	KeyFile string `toml:"key_file"`
}

type workerEnvConfig struct {
	Name          string `toml:"name"`
	Type          string `toml:"type"`
	BroadcastAddr string `toml:"broadcast_addr"`
}

type schedulerConfig struct {
	Addr string `toml:"addr"`
}

type scrapeConfig struct {
	Interval                    string       `toml:"scrape_interval"`
	AttemptServiceConsolidation bool         `toml:"attempt_service_consolidation"`
	Whitelisting                bool         `toml:"whitelisting"`
	ScrapeLabels                scrapeLabels `toml:"scrape_config.labels"`
}

type scrapeLabels struct {
	whitelist string `toml:"whitelist"`
	version   string `toml:"version"`
	name      string `toml:"name"`
}

type tomlConfig struct {
	WorkerEnvConfig workerEnvConfig  `toml:"environment"`
	GRPCServer      grpcServerConfig `toml:"grpc_server"`
	Scheduler       schedulerConfig  `toml:"scheduler"`
	ScrapeConfig    scrapeConfig     `toml:"scrape_config"`
}

var (
	config tomlConfig
)

func loadConfig() {
	var localConfig tomlConfig
	if _, err := toml.DecodeFile("worker/config.toml", &localConfig); err != nil {
		log.Fatalln("Error decoding config file", err)
	}

	config = localConfig
}

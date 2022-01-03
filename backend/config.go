package main

import (
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

//go:embed config_default.toml
var defaultConfigBytes []byte

type Config struct {
	Port    int
	History history
}

type history struct {
	ImageDir string
	Interval duration
	Ttl      duration
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func LoadConfig() (*Config, error) {
	var config Config
	var err error

	// Load default config
	_, err = toml.Decode(string(defaultConfigBytes), &config)
	if err != nil {
		log.Println("Failed to decode default config:", err.Error())
		return nil, err
	}

	// Load user config
	log.Println("Loading config: config.toml")
	userConfigBytes, err := os.ReadFile("config.toml")
	if err != nil {
		log.Println("Failed to load config.toml:", err.Error())
		log.Printf("Config: %+v", config)
		return &config, err
	}
	_, err = toml.Decode(string(userConfigBytes), &config)
	if err != nil {
		log.Println("Failed to decode config.toml:", err.Error())
		log.Printf("Config: %+v", config)
		return &config, err
	}

	log.Printf("Config: %+v", config)
	return &config, nil
}

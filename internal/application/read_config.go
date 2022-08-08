package application

import (
	"flag"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents configuration for the application
type Config struct {
	Logger  Logger  `yaml:"logger"`
	Kafka   Kafka   `yaml:"kafka"`
	Mail    Mail    `yaml:"mail"`
	Compose Compose `yaml:"compose"`
}

// Logger has values for the logger
type Logger struct {
	Level string `yaml:"level"`
}

// Kafka contains parameters for publisher and consumer
type Kafka struct {
	Pub KafkaPub `yaml:"pub"`
	Sub KafkaSub `yaml:"sub"`
}

// KafkaSub keeps values to read messages
type KafkaSub struct {
	Server  string `yaml:"server"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}

// KafkaPub keeps values to send messages
type KafkaPub struct {
	Server string `yaml:"server"`
	Topic  string `yaml:"topic"`
}

// Mail is about to configure smtp messages sender
type Mail struct {
	Rate int `yaml:"rate"`
}

// Compose holds number of workers composing smtp messages simultaneously
type Compose struct {
	Workers int `yaml:"workers"`
}

func getConfig() *Config {
	path := flag.String("c", "./configs/config.yaml", "set path to config yaml-file")
	flag.Parse()

	log.Printf("config file, %s", *path)

	f, err := os.Open(*path)
	if err != nil {
		log.Fatalf("cannot open %s config file: %v", *path, err)
	}
	defer f.Close()

	return readConfigFile(f)
}

// read parses yaml file to get application Config
func readConfigFile(r io.Reader) *Config {

	cfg := &Config{}
	d := yaml.NewDecoder(r)
	if err := d.Decode(cfg); err != nil {
		log.Fatalf("cannot parse config %v", err)
	}
	return cfg
}

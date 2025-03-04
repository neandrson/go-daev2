package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	envErrorStr = "The COMPUTING_POWER environment variable is not set or has an incorrect value."
)

var hostname = flag.String("h", "localhost", "The host name of the orchestrator")
var port = flag.Int("p", 8081, "Port of the orchestrator")

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	flag.Parse()
	if len(*hostname) == 0 {
		*hostname = "localhost"
	}
	if *port <= 0 {
		//fmt.Fprintf(os.Stderr, "Incorrect port: %d\n", *port)
		//os.Exit(1)
		*port = 8081
	}
}

type Config struct {
	ComputingPower int
	Hostname       string
	Port           int
}

func NewConfigFromEnv() (*Config, error) {
	cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || cp <= 0 {
		return nil, fmt.Errorf("%v", envErrorStr)
	}
	cfg := Config{
		ComputingPower: cp,
		Hostname:       *hostname,
		Port:           *port,
	}
	return &cfg, nil
}

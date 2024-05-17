package application

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	Filename            string
	MaxParallelRequests int
}

func NewConfig() *Config {
	return &Config{
		Port:                "8080",
		Filename:            "requests.csv",
		MaxParallelRequests: 5,
	}
}

func (c *Config) LoadFromEnv() {
	port := os.Getenv("PORT")
	if port != "" {
		c.Port = port
	}

	filename := os.Getenv("FILENAME")
	if filename != "" {
		c.Filename = filename
	}

	maxParallelRequests := os.Getenv("MAX_PARALLEL_REQUESTS")
	if maxParallelRequests != "" {
		maxParallelRequestsInt, err := strconv.Atoi(maxParallelRequests)
		if err == nil {
			c.MaxParallelRequests = maxParallelRequestsInt
		}
	}
}

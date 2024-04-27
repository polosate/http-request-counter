package application

import "os"

type Config struct {
	Port     string
	Filename string
}

func NewConfig() *Config {
	return &Config{
		Port:     "8080",
		Filename: "requests.csv",
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
}

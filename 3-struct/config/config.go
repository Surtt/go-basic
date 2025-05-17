package config

import "os"

type Config struct {
	Key string
}

func NewConfig() *Config {
	key := os.Getenv("KEY")
	if key == "" {
		panic("KEY is not set")
	}
	return &Config{
		Key: key,
	}
}

package config

import "github.com/soner3/flora"

type Config struct {
	flora.Component `flora:"constructor=NewConfig"`
	Port            int
	Host            string
}

func NewConfig() Config {
	return Config{
		Port: 8080,
		Host: "localhost",
	}
}

package config

import "github.com/soner3/mint"

type Config struct {
	mint.Component `mint:"constructor=NewConfig"`
	Port           int
	Host           string
}

func NewConfig() Config {
	return Config{
		Port: 8080,
		Host: "localhost",
	}
}

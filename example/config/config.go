package config

import "github.com/soner3/weld"

type Config struct {
	weld.Component `weld:"constructor=NewConfig"`
	Port           int
	Host           string
}

func NewConfig() Config {
	return Config{
		Port: 8080,
		Host: "localhost",
	}
}

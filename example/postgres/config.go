package postgres

import "github.com/soner3/weld"

type Config struct {
	weld.Component
}

func NewConfig() Config {
	return Config{}
}

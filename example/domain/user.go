package domain

import (
	"fmt"

	"github.com/soner3/weld"
	"github.com/soner3/weld/example/config"
)

type UserRepository interface {
	GetUserName() string
}

type UserService struct {
	weld.Component `weld:"constructor=BuildUserService"`
	repo           UserRepository
	cfg            config.Config
}

func BuildUserService(repo UserRepository, cfg config.Config) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *UserService) PrintUser() {
	fmt.Printf("Hello, %s!\n", s.repo.GetUserName())
}

type DummyConfig struct {
	Host string
	Port int
}

func NewDummyConfig() *DummyConfig {
	return &DummyConfig{Host: "localhost", Port: 8080}
}

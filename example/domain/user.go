package domain

import (
	"fmt"

	"github.com/soner3/weld"
)

type UserRepository interface {
	GetUserName() string
}

type UserService struct {
	weld.Component `weld:"constructor=BuildUserService"`
	repo           UserRepository
}

func BuildUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
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

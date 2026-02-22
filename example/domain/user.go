/*
Copyright © 2026 Soner Astan astansoner@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package domain

import (
	"fmt"

	"github.com/soner3/mint"
	"github.com/soner3/mint/example/config"
)

type UserRepository interface {
	GetUserName() string
}

type UserService struct {
	mint.Component `mint:"constructor=BuildUserService"`
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

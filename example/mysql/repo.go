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
package mysql

import (
	"github.com/soner3/flora"
	"github.com/soner3/flora/example/config"
)

type MysqlRepository struct {
	flora.Component
	cfg config.Config
}

func NewMysqlRepository(cfg config.Config) *MysqlRepository {
	return &MysqlRepository{
		cfg: cfg,
	}
}

func (r *MysqlRepository) GetUserName() string {
	return "Foo (from MySQL)"
}

func (r *MysqlRepository) String() string {
	return "MysqlRepository"
}

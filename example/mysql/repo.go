package mysql

import (
	"github.com/soner3/weld"
	"github.com/soner3/weld/example/config"
)

type MysqlRepository struct {
	weld.Component
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

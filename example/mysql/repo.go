package mysql

import (
	"github.com/soner3/weld"
)

type MysqlRepository struct {
	weld.Component
}

func NewMysqlRepository() *MysqlRepository {
	return &MysqlRepository{}
}

func (r *MysqlRepository) GetUserName() string {
	return "Foo (from MySQL)"
}

func (r *MysqlRepository) String() string {
	return "MysqlRepository"
}

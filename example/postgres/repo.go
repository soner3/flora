package postgres

import (
	"github.com/soner3/weld"
)

type PostgresRepository struct {
	weld.Component `weld:"primary"`
}

func NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{}
}

func (r *PostgresRepository) GetUserName() string {
	return "Foo (from Postgres)"
}

func (r *PostgresRepository) String() string {
	return "PostgresRepository"
}

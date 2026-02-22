package postgres

import (
	"github.com/soner3/weld"
)

type PostgresRepository struct {
	weld.Component `weld:"primary"`
	cfg            Config
}

func NewPostgresRepository(cfg Config) *PostgresRepository {
	return &PostgresRepository{
		cfg: cfg,
	}
}

func (r *PostgresRepository) GetUserName() string {
	return "Foo (from Postgres)"
}

func (r *PostgresRepository) String() string {
	return "PostgresRepository"
}

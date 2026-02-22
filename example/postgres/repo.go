package postgres

import (
	"github.com/soner3/weld"
	"github.com/soner3/weld/example/config"
)

type PostgresRepository struct {
	weld.Component `weld:"primary"`
	cfg            config.Config
}

func NewPostgresRepository(cfg config.Config) *PostgresRepository {
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

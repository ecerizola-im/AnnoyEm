package config

import (
	"github.com/ecerizola-im/AnnoyEm/internal/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryConfig struct {
	Type     common.RepositoryType
	Postgres *pgxpool.Pool
}

func (r RepositoryConfig) GetRepoType() common.RepositoryType { return r.Type }
func (r RepositoryConfig) GetPostgresDB() *pgxpool.Pool       { return r.Postgres }

func (r RepositoryConfig) CleanResources() {
	if r.Type == common.TypePostgres {
		r.Postgres.Close()
	}
}

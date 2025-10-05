package config

import (
	"database/sql"

	"github.com/ecerizola-im/AnnoyEm/internal/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryConfig struct {
	Type     common.RepositoryType
	Postgres *pgxpool.Pool
	SQLite   *sql.DB
}

func (r RepositoryConfig) GetRepoType() common.RepositoryType { return r.Type }
func (r RepositoryConfig) GetPostgresDB() *pgxpool.Pool       { return r.Postgres }
func (r RepositoryConfig) GetSQLiteDB() *sql.DB               { return r.SQLite }

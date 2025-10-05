package repository

import (
	"fmt"

	"github.com/ecerizola-im/AnnoyEm/internal/common"
	"github.com/ecerizola-im/AnnoyEm/internal/memes"
	"github.com/ecerizola-im/AnnoyEm/internal/repository/implementation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config interface {
	GetRepoType() common.RepositoryType
	GetPostgresDB() *pgxpool.Pool
}

func NewRepository(cfg Config) (memes.Repository, error) {
	switch cfg.GetRepoType() {
	case common.TypeMemory:
		return implementation.NewMemoryRepository(), nil
	case common.TypePostgres:
		if cfg.GetPostgresDB() == nil {
			return nil, fmt.Errorf("postgres repository requires DB pool")
		}
		return implementation.NewPostgresRepository(cfg.GetPostgresDB()), nil
	default:
		return nil, fmt.Errorf("unknown repository type %q", cfg.GetRepoType())
	}
}

package repository

import (
	"fmt"

	"github.com/ecerizola-im/AnnoyEm/internal/memes"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Type string

const (
	TypeMemory   Type = "memory"
	TypePostgres Type = "postgres"
)

type Config struct {
	Type Type
	DB   *pgxpool.Pool
}

func NewRepository(cfg Config) (memes.Repository, error) {
	switch cfg.Type {
	case TypeMemory:
		return NewMemoryRepository(), nil
	case TypePostgres:
		if cfg.DB == nil {
			return nil, fmt.Errorf("postgres repository requires DB pool")
		}
		return NewPostgresRepository(cfg.DB), nil
	default:
		return nil, fmt.Errorf("unknown repository type %q", cfg.Type)
	}
}

package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, data io.Reader) (string, error)
	Delete(ctx context.Context, fileName string) error
	Get(ctx context.Context, fileName string) (io.ReadCloser, error)
}

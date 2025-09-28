package memes

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, r *Meme) (int64, error)
	FindByID(id int64) (*Meme, error)
	List() ([]Meme, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, r *Meme) error
}

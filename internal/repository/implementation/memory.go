package implementation

import (
	"context"
	"fmt"

	"github.com/ecerizola-im/AnnoyEm/internal/memes"
)

type MemoryRepository struct {
	data  map[int64]*memes.Meme
	maxID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data:  make(map[int64]*memes.Meme),
		maxID: 0,
	}
}

func (r *MemoryRepository) Create(ctx context.Context, meme *memes.Meme) (int64, error) {

	newMemeId := r.maxID + 1

	memeCopy := *meme
	memeCopy.ID = newMemeId
	r.data[newMemeId] = &memeCopy
	r.maxID = newMemeId

	return newMemeId, nil
}

func (r *MemoryRepository) FindByID(id int64) (*memes.Meme, error) {
	meme, exists := r.data[id]
	if !exists {
		return nil, fmt.Errorf("meme not found")
	}
	return meme, nil
}

func (r *MemoryRepository) List() ([]memes.Meme, error) {
	memesList := make([]memes.Meme, 0, len(r.data))
	for _, meme := range r.data {
		memesList = append(memesList, *meme)
	}
	return memesList, nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id int64) error {
	_, exists := r.data[id]
	if !exists {
		return fmt.Errorf("meme not found")
	}
	delete(r.data, id)
	return nil
}

func (r *MemoryRepository) Update(ctx context.Context, meme *memes.Meme) error {
	_, exists := r.data[meme.ID]
	if !exists {
		return fmt.Errorf("meme not found")
	}
	r.data[meme.ID] = meme
	return nil
}

func (r *MemoryRepository) Cleanup() {

}

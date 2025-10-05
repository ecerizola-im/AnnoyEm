package memes

import (
	"context"
	"fmt"
	"io"
	"time"

	memeModel "github.com/ecerizola-im/AnnoyEm/internal/model/meme"
	"github.com/ecerizola-im/AnnoyEm/internal/storage"
)

type MemeService struct {
	repository Repository
	storage    storage.Storage
	now        func() time.Time
}

func NewMemeService(repository Repository, storage storage.Storage) *MemeService {
	return &MemeService{
		repository: repository,
		storage:    storage,
		now:        time.Now,
	}
}

func (svc *MemeService) AddMeme(ctx context.Context, fileName string, data io.Reader) (int64, error) {

	meme := GetNewMeme(fileName)

	memeId, err := svc.repository.Create(ctx, meme)

	if err != nil {
		return 0, fmt.Errorf("failed to create meme: %w", err)
	}

	meme.ID = memeId

	uploadedFileId, err := svc.storage.Save(ctx, data)

	if err != nil {
		meme.Status = memeModel.Failed
		meme.UpdatedAt = svc.now()

		if err = svc.repository.Update(ctx, meme); err != nil {
			return 0, fmt.Errorf("failed to update meme to status %d for memeID %d: %w", meme.Status, meme.ID, err)
		}

		return 0, fmt.Errorf("failed to save file for memeID %d: %w", meme.ID, err)
	}

	meme.UUID = &uploadedFileId

	currentTime := svc.now()
	meme.Status = memeModel.Processed
	meme.ProcessedAt = &currentTime
	meme.UpdatedAt = currentTime

	if err = svc.repository.Update(ctx, meme); err != nil {
		return 0, fmt.Errorf("failed to update meme to status %d for memeID %d: %w", meme.Status, meme.ID, err)
	}

	return memeId, nil

}

func (svc *MemeService) GetMemeFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	return svc.storage.Get(ctx, fileID)
}

func (svc *MemeService) GetMeme(ctx context.Context, memeId int64) (*Meme, error) {
	meme, err := svc.repository.FindByID(memeId)
	if err != nil {
		return nil, fmt.Errorf("failed to find meme: %w", err)
	}

	return meme, nil
}

func (svc *MemeService) GetMemes(ctx context.Context) ([]Meme, error) {
	memes, err := svc.repository.List()

	if err != nil {
		return nil, fmt.Errorf("failed to get memes: %w", err)
	}
	return memes, nil
}

func GetNewMeme(fileName string) *Meme {

	currentTime := time.Now()

	return &Meme{
		ID:               0,
		UserID:           0,
		OriginalFileName: fileName,
		MimeType:         "",
		SizeBytes:        0,
		UUID:             nil,
		Status:           memeModel.Pending,
		Category:         "",
		CreatedAt:        currentTime,
		ProcessedAt:      nil,
		UpdatedAt:        currentTime,
	}
}

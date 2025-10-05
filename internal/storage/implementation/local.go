package implementation

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	os.MkdirAll(basePath, os.ModePerm)
	return &LocalStorage{BasePath: basePath}
}

func (s *LocalStorage) Save(ctx context.Context, data io.Reader) (string, error) {

	fileId := uuid.NewString()
	filePath := s.BasePath + "/" + fileId
	newFile, err := os.Create(filePath)

	if err != nil {
		return "", err
	}

	defer newFile.Close()

	_, err = io.Copy(newFile, data)
	if err != nil {
		return "", err
	}

	return fileId, nil
}

func (s *LocalStorage) Delete(ctx context.Context, fileName string) error {

	path := filepath.Join(s.BasePath, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Get(ctx context.Context, fileName string) (io.ReadCloser, error) {

	path := filepath.Join(s.BasePath, fileName)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

package storage

import (
	"fmt"

	"github.com/ecerizola-im/AnnoyEm/internal/common"
	"github.com/ecerizola-im/AnnoyEm/internal/storage/implementation"
)

type StorageConfig interface {
	GetStorageType() common.StorageType
	GetContainerName() string
	GetLocalStorageBasePath() string
}

func CreateStorage(cfg StorageConfig) (Storage, error) {
	switch cfg.GetStorageType() {
	case common.AzureBlob:
		storeConfig := implementation.AzureBlobConfig{ContainerName: cfg.GetContainerName()}
		storage, err := implementation.NewAzureBlobStorage(storeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure Blob storage: %v", err)
		}
		return storage, nil
	case common.LocalStorage:
		return implementation.NewLocalStorage(cfg.GetLocalStorageBasePath()), nil
	default:
		return nil, fmt.Errorf("unknown storage type %q", cfg.GetStorageType())
	}
}

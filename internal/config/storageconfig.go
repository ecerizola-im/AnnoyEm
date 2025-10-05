package config

import "github.com/ecerizola-im/AnnoyEm/internal/common"

type StorageConfig struct {
	Type                 common.StorageType
	ContainerName        string
	LocalStorageBasePath string
}

func (s StorageConfig) GetStorageType() common.StorageType {
	return s.Type
}

func (s StorageConfig) GetContainerName() string {
	return s.ContainerName
}

func (s StorageConfig) GetLocalStorageBasePath() string {
	return s.LocalStorageBasePath
}

func GetStorageConfig() StorageConfig {
	return StorageConfig{
		Type:                 common.AzureBlob,
		ContainerName:        "memes",
		LocalStorageBasePath: "./data/receipts",
	}
}

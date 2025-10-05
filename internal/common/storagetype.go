package common

type StorageType string

const (
	AzureBlob    StorageType = "azure_blob"
	LocalStorage StorageType = "local_storage"
)

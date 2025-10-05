package implementation

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

type AzureBlobConfig struct {
	ContainerName string
}

type AzureBlobStorage struct {
	client    *azblob.Client
	container string
}

func handleError(err error) {
	if err != nil {
		log.Printf("Azure Blob Storage error: %v", err)
	}
}

func NewAzureBlobStorage(config AzureBlobConfig) (*AzureBlobStorage, error) {

	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)

	if err != nil {
		log.Fatalf("failed to create Azure identity: %v", err)
		return nil, err
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	client, err := azblob.NewClientWithSharedKeyCredential(url, credential, nil)

	if err != nil {
		log.Fatalf("failed to create Azure Blob client: %v", err)
		return nil, err
	}

	return &AzureBlobStorage{
		client:    client,
		container: config.ContainerName,
	}, nil
}

func (s *AzureBlobStorage) Save(ctx context.Context, data io.Reader) (string, error) {
	// Upload to data to blob storage

	uuid := uuid.NewString()
	_, err := s.client.UploadStream(ctx, s.container, uuid, data, &azblob.UploadStreamOptions{})

	if err != nil {
		handleError(err)
		return "", err
	}

	return uuid, nil
}

func (s *AzureBlobStorage) Delete(ctx context.Context, fileName string) error {

	_, err := s.client.DeleteBlob(ctx, s.container, fileName, nil)
	handleError(err)
	return err
}

func (s *AzureBlobStorage) Get(ctx context.Context, fileName string) (io.ReadCloser, error) {

	contentLength := int64(0)

	dr, err := s.client.DownloadStream(ctx, s.container, fileName, &azblob.DownloadStreamOptions{})

	if err != nil {
		handleError(err)
		return nil, err
	}

	resultStream := dr.Body

	stream := streaming.NewResponseProgress(
		resultStream,
		func(bytesTransferred int64) {
			fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)
		},
	)

	return stream, nil
}

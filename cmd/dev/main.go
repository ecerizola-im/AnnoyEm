// cmd/AnnoyEm/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ecerizola-im/AnnoyEm/internal/app"
	"github.com/ecerizola-im/AnnoyEm/internal/config"
	"github.com/ecerizola-im/AnnoyEm/internal/memes"
	"github.com/ecerizola-im/AnnoyEm/internal/memes/repository"
	"github.com/ecerizola-im/AnnoyEm/internal/memes/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	config := config.Load()
	//DB
	repoConfig := getRepoConfig(config)

	if config.RepoType == repository.TypePostgres {
		defer repoConfig.DB.Close()
	}

	repo, err := repository.NewRepository(repoConfig)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	//uploadsDir := configureStorageSettings()

	storeConfig := storage.AzureBlobConfig{ContainerName: "memes"}

	store, err := storage.NewAzureBlobStorage(storeConfig)
	if err != nil {
		log.Fatalf("failed to create Azure Blob storage: %v", err)
	}
	// store := storage.NewLocalStorage(uploadsDir)

	// Service
	svc := memes.NewMemeService(repo, store)

	// HTTP
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static"))))

	h := memes.NewHandler(svc)
	h.Register(mux) // exposes POST /memes (from your handler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", app.LoggingMiddleware(mux)))
}

func getRepoConfig(c config.Config) repository.Config {

	repoConfig := repository.Config{
		Type: c.RepoType,
		DB:   configureDbPool(c),
	}
	return repoConfig
}

func configureDbPool(c config.Config) *pgxpool.Pool {

	connectionString := c.ConnectionString
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return dbpool
}

func configureStorageSettings() string {
	// Where to store uploaded files locally
	uploadsDir := filepath.Join(".", "data", "receipts")
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		log.Fatalf("failed to create uploads dir: %v", err)
	}
	return uploadsDir
}

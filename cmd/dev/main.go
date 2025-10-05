// cmd/AnnoyEm/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ecerizola-im/AnnoyEm/internal/app"
	"github.com/ecerizola-im/AnnoyEm/internal/common"
	"github.com/ecerizola-im/AnnoyEm/internal/config"
	"github.com/ecerizola-im/AnnoyEm/internal/memes"
	"github.com/ecerizola-im/AnnoyEm/internal/repository"
	"github.com/ecerizola-im/AnnoyEm/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "modernc.org/sqlite"
)

func main() {

	config := config.Load()
	//DB
	repoConfig, err := getRepoConfig(config)

	if err != nil {
		log.Fatalf("failed to get repository config: %v", err)
	}

	repo, err := repository.NewRepository(repoConfig)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	defer repo.Cleanup()

	//uploadsDir := configureStorageSettings()

	store, err := storage.CreateStorage(config.Storage)
	if err != nil {
		log.Fatalf("failed to create local storage: %v", err)
	}

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

func getRepoConfig(c config.Config) (*config.RepositoryConfig, error) {

	switch c.RepoType {
	case common.TypeMemory:
		return &config.RepositoryConfig{Type: common.TypeMemory}, nil
	case common.TypePostgres:
		return &config.RepositoryConfig{Type: common.TypePostgres, Postgres: configureDbPool(c)}, nil
	case common.TypeSQLite:
		return &config.RepositoryConfig{Type: common.TypeSQLite, SQLite: configureSQLiteDB(c)}, nil
	default:
		return nil, fmt.Errorf("unknown repository type: %v", c.RepoType)
	}
}

func configureDbPool(c config.Config) *pgxpool.Pool {

	connectionString := c.ConnectionString
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return dbpool
}

func configureSQLiteDB(c config.Config) *sql.DB {
	db, err := sql.Open("sqlite", c.EmbeddedDatabase)
	if err != nil {
		log.Fatalf("Unable to create SQLite connection: %v\n", err)
	}
	return db
}

package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ecerizola-im/AnnoyEm/internal/common"
)

type DatabaseConfig struct {
	URL      string
	Port     string
	UserName string
	Password string
	DBName   string
}

type Config struct {
	Port             string
	ConnectionString string
	RepoType         common.RepositoryType
	MaxUploadBytes   int64
	UploadsDir       string
	Storage          StorageConfig
}

func Load() Config {
	c := Config{
		Port:             "8080",
		ConnectionString: LoadDatabaseConfig(),
		RepoType:         common.TypePostgres,
		MaxUploadBytes:   10 << 20,
		UploadsDir:       filepath.Join(".", "data", "receipts"),
		Storage: StorageConfig{
			Type:                 "azure_blob",
			ContainerName:        "memes",
			LocalStorageBasePath: filepath.Join(".", "data", "uploads"),
		},
	}

	if v := os.Getenv("AnnoyEm_PORT"); v != "" {
		c.Port = v
	}

	if v := os.Getenv("AnnoyEm_REPO_TYPE"); v != "" {
		switch strings.ToLower(v) {
		case "memory":
			c.RepoType = common.TypeMemory
		case "postgres":
			c.RepoType = common.TypePostgres
		}
	}

	if v := os.Getenv("AnnoyEm_MAX_UPLOAD_BYTES"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			c.MaxUploadBytes = n
		}
	}
	if v := os.Getenv("AnnoyEm_UPLOADS_DIR"); v != "" {
		c.UploadsDir = v
	}
	return c
}

func LoadDatabaseConfig() string {

	c := DatabaseConfig{
		URL:      "localhost",
		Port:     "5432",
		UserName: "annoyem",
		Password: "annoyem",
		DBName:   "annoyem",
	}

	if v := os.Getenv("AnnoyEm_DATABASE_URL"); v != "" {
		c.URL = v
	}
	if v := os.Getenv("AnnoyEm_DATABASE_PORT"); v != "" {
		c.Port = v
	}
	if v := os.Getenv("AnnoyEm_DATABASE_USERNAME"); v != "" {
		c.UserName = v
	}
	if v := os.Getenv("AnnoyEm_DATABASE_PASSWORD"); v != "" {
		c.Password = v
	}
	if v := os.Getenv("AnnoyEm_DATABASE_NAME"); v != "" {
		c.DBName = v
	}

	var connString = "postgres://" + c.UserName + ":" + c.Password + "@" + c.URL + ":" + c.Port + "/" + c.DBName + "?sslmode=disable"

	return connString
}

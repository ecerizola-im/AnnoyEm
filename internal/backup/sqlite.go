// internal/backup/sqlite.go
package backup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ecerizola-im/AnnoyEm/internal/storage"
)

// BackupSQLiteUsingStorage creates a hot snapshot (VACUUM INTO),
// then uploads it via your Storage. Returns the storage key/id.
func BackupSQLiteUsingStorage(ctx context.Context, db *sql.DB, dbPath string, storage storage.Storage) (string, error) {

	if _, err := db.ExecContext(ctx, `PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		return "", fmt.Errorf("wal checkpoint: %w", err)
	}

	ts := time.Now().UTC().Format("20060102-150405")
	tmpDir := filepath.Dir(dbPath)
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("backup-%s.db", ts))

	if _, err := db.ExecContext(ctx, `VACUUM INTO ?;`, tmpFile); err != nil {
		return "", fmt.Errorf("vacuum into: %w", err)
	}
	defer os.Remove(tmpFile)

	// 3) stream it to storage
	f, err := os.Open(tmpFile)
	if err != nil {
		return "", fmt.Errorf("open temp backup: %w", err)
	}
	defer f.Close()

	// prefer a stable path if the storage supports naming
	storageFileName := filepath.Join("backups", "sqlite",
		time.Now().UTC().Format("2006"),
		time.Now().UTC().Format("01"),
		fmt.Sprintf("app-%s.db", ts),
	)

	result, err := storage.Save(ctx, f, storageFileName)

	if err != nil {
		return "", fmt.Errorf("upload backup to storage: %w", err)
	}

	return result, nil
}

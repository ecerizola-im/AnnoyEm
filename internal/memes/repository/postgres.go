package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ecerizola-im/AnnoyEm/internal/memes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, meme *memes.Meme) (int64, error) {
	const q = `
        INSERT INTO memes.meme (user_id, original_file_name, mime_type, size_bytes, uuid, upload_status_id,
			   category, created_at, processed_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id`

	now := time.Now().UTC()
	meme.CreatedAt = now
	meme.UpdatedAt = now

	var id int64
	err := r.db.QueryRow(ctx, q,
		meme.UserID,
		meme.OriginalFileName,
		meme.MimeType,
		meme.SizeBytes,
		meme.UUID,
		meme.Status,
		meme.Category,
		meme.CreatedAt,
		meme.ProcessedAt,
		meme.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert receipt: %w", err)
	}
	return id, nil
}

func (r *PostgresRepository) FindByID(id int64) (*memes.Meme, error) {
	const q = `
        SELECT id, user_id, original_file_name, mime_type, size_bytes, uuid, upload_status_id,
			   category, created_at, processed_at, updated_at
        FROM memes.meme
        WHERE id = $1`
	row := r.db.QueryRow(context.Background(), q, id)

	var rec memes.Meme
	err := row.Scan(
		&rec.ID,
		&rec.UserID,
		&rec.OriginalFileName,
		&rec.MimeType,
		&rec.SizeBytes,
		&rec.UUID,
		&rec.Status,
		&rec.Category,
		&rec.CreatedAt,
		&rec.ProcessedAt,
		&rec.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find receipt: %w", err)
	}
	return &rec, nil
}

func (r *PostgresRepository) List() ([]memes.Meme, error) {
	const q = `
        SELECT id, user_id, original_file_name, mime_type, size_bytes, uuid, upload_status_id,
			   category, created_at, processed_at, updated_at
        FROM memes.meme
        ORDER BY created_at DESC`
	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, fmt.Errorf("list meme: %w", err)
	}
	defer rows.Close()

	var out []memes.Meme
	for rows.Next() {
		var rec memes.Meme
		if err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.OriginalFileName,
			&rec.MimeType,
			&rec.SizeBytes,
			&rec.UUID,
			&rec.Status,
			&rec.Category,
			&rec.CreatedAt,
			&rec.ProcessedAt,
			&rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan meme: %w", err)
		}
		out = append(out, rec)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return out, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM receipts.receipt WHERE id = $1`
	if _, err := r.db.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete meme: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, rec *memes.Meme) error {
	rec.UpdatedAt = nowUTC()
	const q = `
        UPDATE memes.meme
        SET user_id = $1, original_file_name = $2, mime_type = $3, size_bytes = $4, uuid = $5, upload_status_id = $6,
			category = $7, created_at = $8, processed_at = $9, updated_at = $10
        WHERE id = $11`
	_, err := r.db.Exec(ctx, q,
		rec.UserID,
		rec.OriginalFileName,
		rec.MimeType,
		rec.SizeBytes,
		rec.UUID,
		rec.Status,
		rec.Category,
		rec.CreatedAt,
		rec.ProcessedAt,
		rec.UpdatedAt,
		rec.ID,
	)
	if err != nil {
		return fmt.Errorf("update receipt: %w", err)
	}
	return nil
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

package meme

import (
	"time"
)

type Meme struct {
	ID               int64            `json:"id"`
	UserID           int64            `json:"user_id"`
	OriginalFileName string           `json:"original_file_name"`
	MimeType         string           `json:"mime_type"`
	SizeBytes        int64            `json:"size_bytes"`
	UUID             *string          `json:"uuid"`
	Status           FileUploadStatus `json:"status"`
	Category         string           `json:"category"`
	CreatedAt        time.Time        `json:"created_at"`
	ProcessedAt      *time.Time       `json:"processed_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type FileMetaData struct {
	SizeBytes int64
	MimeType  string
}

type FileUploadStatus int8

const (
	Unspecified FileUploadStatus = iota
	Pending
	Processed
	Failed
)

func (s FileUploadStatus) String() string {
	switch s {
	case Pending:
		return "pending"
	case Processed:
		return "processed"
	case Failed:
		return "failed"
	default:
		return "unknown"
	}
}

type PaymentStatus int8

const (
	Unpaid PaymentStatus = iota
	Paid
)

func (p PaymentStatus) String() string {
	switch p {
	case Unpaid:
		return "unpaid"
	case Paid:
		return "paid"
	default:
		return "unknown"
	}
}

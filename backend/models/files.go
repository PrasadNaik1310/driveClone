package models

import "time"

type File struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	FolderID   *string   `json:"folder_id"` // nullable (root files)
	OwnerID    uint      `json:"owner_id"`
	StorageKey string    `json:"storage_key"` // S3 key
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

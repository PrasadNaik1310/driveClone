package models

import "time"

type File struct {
	ID         string    `gorm:"primaryKey" json:"fileid"`
	Name       string    `json:"name"`
	FolderID   *string   `json:"folderid"` // nullable (root files)
	OwnerID    uint      `json:"ownerid"`
	StorageKey string    `json:"storagekey"` // S3 key
	Size       int64     `json:"size"`
	MimeType   string    `json:"mimetype"`
	CreatedAt  time.Time `json:"createdat"`
	UpdatedAt  time.Time `json:"updatedat"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)
type  Folder struct{
gorm.Model
FolderName string `json:"foldername"`
OwnerId uint `json:"ownerid"`
ParentId string `json:"parentid"`
CreatedAt time.Time  `json:"createdat"`
UpdatedAt time.Time `json:"updatedat"`
Uuid string 	`json:"uuid"`
FolderId  string `json:"folderid" gorm:"primarykey"`



}

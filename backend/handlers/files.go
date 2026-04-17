package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/PrasadNaik1310/driveClone/models"
	"github.com/PrasadNaik1310/driveClone/db"
)

func UploadFile(c *gin.Context) {
	userIdRaw, exists := c.Get("user_id")
	if ! exists{
	c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthenticated user "})
	log.Printf("Unauthenticated access attempt from %s", c.ClientIP())
	return
	}
UserId, ok := userIdRaw.(uint)
if !ok {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
	log.Printf("Invalid user_id type: expected uint, got %T", userIdRaw)
	return}
	newFileName := c.PostForm("newfilename")

	var fileHeaderName string
	var fileErr error
	var file = new(struct {
		Filename string
	})

	for _, key := range []string{"file", "files", "upload", "document"} {
		fh, err := c.FormFile(key)
		if err == nil {
			fileHeaderName = key
			file.Filename = fh.Filename
			fileErr = nil
			break
		}
		fileErr = err
	}

	if fileHeaderName == "" {
		form, _ := c.MultipartForm()
		receivedKeys := make([]string, 0)
		if form != nil {
			for k := range form.Value {
				receivedKeys = append(receivedKeys, k)
			}
			for k := range form.File {
				receivedKeys = append(receivedKeys, k)
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":                "File not found in request. Send multipart/form-data with a file field named 'file'",
			"received_content_type": c.ContentType(),
			"received_form_keys":    receivedKeys,
			"tried_file_fields":     []string{"file", "files", "upload", "document"},
			"details":              fileErr.Error(),
		})
		return
	}

	fh, err := c.FormFile(fileHeaderName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read uploaded file"})
		return
	}
	folderId := c.PostForm("folder_id")
	// creating a storage for each file because i cant store raw files as there as chances of file duplicate names .
	ext := filepath.Ext(fh.Filename)
	fileID := uuid.New().String()
	storageName := fileID + ext

	os.MkdirAll("uploads", os.ModePerm)
	filePath := filepath.Join("uploads", storageName)
	err = c.SaveUploadedFile(fh, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Unable to save file"})
		return
	}
	log.Printf("Save file successfull at %s", filePath)

	if strings.TrimSpace(newFileName) == "" {
		newFileName = fh.Filename
	}

	log.Printf("Got file :-> %s (field: %s) from client %s ", storageName, fileHeaderName, c.ClientIP())
	/*c.JSON(http.StatusOK, gin.H{
		"file":        fh.Filename,
		"newfilename": newFileName,
		"storageName": storageName,
		"folder_id":   folderId,
		"file_field":  fileHeaderName,
	})*/
	log.Printf("Moving for db migration for file %s (field: %s) from client %s ", storageName, fileHeaderName, c.ClientIP())
	
	var fileModel models.File 

fileModel.ID = fileID
fileModel.Name = newFileName
fileModel.OwnerID  = UserId
fileModel.StorageKey = filePath
fileModel.FolderID = &folderId
fileModel.Size = fh.Size
fileModel.MimeType = fh.Header.Get("Content-Type")

if err := db.DB.Create(&fileModel).Error; err != nil {
	log.Printf("Error migrating file %s to DB :%v",fileModel.Name, err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"Error": "Could not migrate the file from disk to DB . Try again ",
	})
	log.Printf("Deleting file %s from disk", fileModel.Name)
deleteErr := os.Remove(filePath)
if deleteErr != nil {
	log.Printf("Error deleting file %s from disk after failed migration", fileModel.StorageKey, deleteErr)
}
	return
}
log.Printf("File %s migrated to DB , DONEEEE",fileModel.StorageKey)
c.JSON(http.StatusOK, gin.H{
	"message": "File uploaded and migrated to DB successfully",
	"file":    fileModel,
	})

}
func DownloadFile(c *gin.Context) {

	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file id required"})
		return
	}

	// 🔹 Get user
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	userID := userIDRaw.(uint)

	// 🔹 Fetch file from DB
	var file models.File
	if err := db.DB.Where("id = ? AND owner_id = ?", fileID, userID).
		First(&file).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// 🔹 Serve file
	c.FileAttachment(file.StorageKey, file.Name)
}

